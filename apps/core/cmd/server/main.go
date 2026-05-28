package main

import (
	"context"
	stdhttp "net/http"
	"os"
	"os/signal"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humaecho"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/fanboykun/smokery/apps/core/internal/adapter/inproc"
	minioadapter "github.com/fanboykun/smokery/apps/core/internal/adapter/minio"
	"github.com/fanboykun/smokery/apps/core/internal/adapter/postgres"
	"github.com/fanboykun/smokery/apps/core/internal/app"
	"github.com/fanboykun/smokery/apps/core/internal/config"
	deliveryhttp "github.com/fanboykun/smokery/apps/core/internal/delivery/http"
	"github.com/fanboykun/smokery/apps/core/internal/frontend"
	"github.com/fanboykun/smokery/apps/core/internal/runner"
)

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	cfg, err := config.Load()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to load config")
	}

	// --- Adapters ---
	pool, err := pgxpool.New(context.Background(), cfg.DatabaseURL)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to connect to database")
	}
	defer pool.Close()

	projectRepo := postgres.NewProjectRepo(pool)
	specRepo := postgres.NewSpecRepo(pool)
	operationRepo := postgres.NewOperationRepo(pool)
	runRepo := postgres.NewRunRepo(pool)
	commentRepo := postgres.NewCommentRepo(pool)
	artifactRepo := postgres.NewArtifactRepo(pool)

	eventBus := inproc.NewEventBus()
	rnr := runner.New(runner.DefaultOptions())
	worker := inproc.NewWorker(runRepo, eventBus, rnr)

	// Set up blob store for artifact persistence
	blobStore, err := minioadapter.New(minioadapter.Config{
		Endpoint:  cfg.MinioEndpoint,
		AccessKey: cfg.MinioAccessKey,
		SecretKey: cfg.MinioSecretKey,
		Bucket:    "smokery",
		UseSSL:    false,
	})
	if err != nil {
		log.Warn().Err(err).Msg("MinIO unavailable, artifacts will not be persisted")
	} else {
		worker.WithArtifacts(blobStore, artifactRepo)
	}

	// --- App services ---
	projectSvc := app.NewProjectService(projectRepo)
	specSvc := app.NewSpecService(specRepo, operationRepo)
	operationSvc := app.NewOperationService(operationRepo)
	runSvc := app.NewRunService(runRepo, worker)
	reportSvc := app.NewReportService(runRepo)
	commentSvc := app.NewCommentService(commentRepo)
	artifactSvc := app.NewArtifactService(artifactRepo)
	planSvc := app.NewPlanService(specRepo, operationRepo)

	// --- Echo + Huma ---
	e := echo.New()
	e.HideBanner = true
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{stdhttp.MethodGet, stdhttp.MethodPost, stdhttp.MethodPut, stdhttp.MethodDelete, stdhttp.MethodOptions},
	}))

	humaConfig := huma.DefaultConfig("Smokery API", "1.0.0")
	humaConfig.Servers = []*huma.Server{{URL: "http://localhost:" + cfg.Port}}
	api := humaecho.New(e, humaConfig)

	// --- Delivery ---
	healthChecker := &serverHealthChecker{pool: pool, blob: blobStore}
	deliveryhttp.RegisterHealthCheck(api, healthChecker)
	deliveryhttp.RegisterProjects(api, projectSvc)
	deliveryhttp.RegisterSpecs(api, specSvc)
	deliveryhttp.RegisterOperations(api, operationSvc)
	deliveryhttp.RegisterRuns(api, runSvc)
	deliveryhttp.RegisterReports(api, reportSvc)
	deliveryhttp.RegisterComments(api, commentSvc)
	deliveryhttp.RegisterArtifacts(api, artifactSvc)
	deliveryhttp.RegisterWebSocket(e, eventBus)
	deliveryhttp.RegisterPlan(api, planSvc)

	// --- Embedded frontend (production only, built with -tags embed_frontend) ---
	if frontendFS := frontend.FS(); frontendFS != nil {
		fsHandler := stdhttp.FileServer(stdhttp.FS(frontendFS))
		e.GET("/*", echo.WrapHandler(stdhttp.StripPrefix("/", fsHandler)))
		log.Info().Msg("serving embedded frontend")
	}

	// --- Lifecycle ---
	retentionTTL := time.Duration(cfg.RetentionDays) * 24 * time.Hour
	cleaner := inproc.NewRetentionCleaner(runRepo, artifactRepo, blobStore, retentionTTL, 1*time.Hour)
	cleaner.Start()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()
	go func() {
		if err := e.Start(":" + cfg.Port); err != nil && err != stdhttp.ErrServerClosed {
			log.Fatal().Err(err).Msg("server error")
		}
	}()
	log.Info().Str("port", cfg.Port).Msg("server started")
	log.Info().Msg("OpenAPI spec at /openapi.json")
	<-ctx.Done()
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	cleaner.Stop()
	_ = e.Shutdown(shutdownCtx)
}

// serverHealthChecker probes PostgreSQL and MinIO connectivity.
type serverHealthChecker struct {
	pool *pgxpool.Pool
	blob *minioadapter.BlobStore
}

func (h *serverHealthChecker) PingDB(ctx context.Context) error {
	return h.pool.Ping(ctx)
}

func (h *serverHealthChecker) PingBlob(ctx context.Context) error {
	if h.blob == nil {
		return nil
	}
	return h.blob.Ping(ctx)
}
