package main

import (
	"context"
	stdhttp "net/http"
	"os"
	"os/signal"
	"path/filepath"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humaecho"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/fanboykun/smokery/apps/core/internal/adapter/fs"
	"github.com/fanboykun/smokery/apps/core/internal/adapter/inproc"
	minioadapter "github.com/fanboykun/smokery/apps/core/internal/adapter/minio"
	"github.com/fanboykun/smokery/apps/core/internal/adapter/postgres"
	sqliteadapter "github.com/fanboykun/smokery/apps/core/internal/adapter/sqlite"
	"github.com/fanboykun/smokery/apps/core/internal/app"
	"github.com/fanboykun/smokery/apps/core/internal/config"
	deliveryhttp "github.com/fanboykun/smokery/apps/core/internal/delivery/http"
	"github.com/fanboykun/smokery/apps/core/internal/frontend"
	"github.com/fanboykun/smokery/apps/core/internal/port"
	"github.com/fanboykun/smokery/apps/core/internal/runner"
)

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	cfg, err := config.Load()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to load config")
	}

	// --- Adapters: DB ---
	var (
		projectRepo   port.ProjectRepo
		specRepo      port.SpecRepo
		operationRepo port.OperationRepo
		runRepo       port.RunRepo
		commentRepo   port.CommentRepo
		artifactRepo  port.ArtifactRepo
		dbPinger      func(ctx context.Context) error
		dbCloser      func()
	)

	switch cfg.DBAdapter {
	case "postgres":
		pool, err := pgxpool.New(context.Background(), cfg.DatabaseURL)
		if err != nil {
			log.Fatal().Err(err).Msg("failed to connect to postgres")
		}
		dbCloser = pool.Close
		dbPinger = pool.Ping
		projectRepo = postgres.NewProjectRepo(pool)
		specRepo = postgres.NewSpecRepo(pool)
		operationRepo = postgres.NewOperationRepo(pool)
		runRepo = postgres.NewRunRepo(pool)
		commentRepo = postgres.NewCommentRepo(pool)
		artifactRepo = postgres.NewArtifactRepo(pool)
		log.Info().Str("adapter", "postgres").Msg("database connected")

	default: // sqlite
		if err := os.MkdirAll(filepath.Dir(cfg.SQLitePath), 0o755); err != nil {
			log.Fatal().Err(err).Msg("failed to create sqlite directory")
		}
		db, err := sqliteadapter.Open(cfg.SQLitePath)
		if err != nil {
			log.Fatal().Err(err).Msg("failed to open sqlite database")
		}
		dbCloser = func() { db.Close() }
		dbPinger = db.Ping
		projectRepo = sqliteadapter.NewProjectRepo(db)
		specRepo = sqliteadapter.NewSpecRepo(db)
		operationRepo = sqliteadapter.NewOperationRepo(db)
		runRepo = sqliteadapter.NewRunRepo(db)
		commentRepo = sqliteadapter.NewCommentRepo(db)
		artifactRepo = sqliteadapter.NewArtifactRepo(db)
		log.Info().Str("adapter", "sqlite").Str("path", cfg.SQLitePath).Msg("database connected")
	}
	defer dbCloser()

	// --- Adapters: Blob storage ---
	var (
		blobStore  port.BlobStore
		blobPinger func(ctx context.Context) error
	)

	switch cfg.StorageAdapter {
	case "minio":
		blob, err := minioadapter.New(minioadapter.Config{
			Endpoint:  cfg.MinioEndpoint,
			AccessKey: cfg.MinioAccessKey,
			SecretKey: cfg.MinioSecretKey,
			Bucket:    "smokery",
			UseSSL:    false,
		})
		if err != nil {
			log.Warn().Err(err).Msg("MinIO unavailable, artifacts will not be persisted")
		} else {
			blobStore = blob
			blobPinger = blob.Ping
			log.Info().Str("adapter", "minio").Msg("blob storage connected")
		}

	default: // fs
		if err := os.MkdirAll(cfg.StoragePath, 0o755); err != nil {
			log.Fatal().Err(err).Msg("failed to create storage directory")
		}
		blob, err := fs.New(cfg.StoragePath)
		if err != nil {
			log.Fatal().Err(err).Msg("failed to create fs blob store")
		}
		blobStore = blob
		blobPinger = func(ctx context.Context) error { return nil }
		log.Info().Str("adapter", "fs").Str("path", cfg.StoragePath).Msg("blob storage configured")
	}

	// --- Inproc services ---
	eventBus := inproc.NewEventBus()
	rnr := runner.New(runner.DefaultOptions())
	worker := inproc.NewWorker(runRepo, eventBus, rnr)
	if blobStore != nil {
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
	healthChecker := &serverHealthChecker{pingDB: dbPinger, pingBlob: blobPinger}
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
	log.Info().Str("port", cfg.Port).Str("db", cfg.DBAdapter).Str("storage", cfg.StorageAdapter).Msg("server started")
	log.Info().Msg("OpenAPI spec at /openapi.json")
	<-ctx.Done()
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	cleaner.Stop()
	_ = e.Shutdown(shutdownCtx)
}

// serverHealthChecker probes DB and blob storage connectivity.
type serverHealthChecker struct {
	pingDB   func(ctx context.Context) error
	pingBlob func(ctx context.Context) error
}

func (h *serverHealthChecker) PingDB(ctx context.Context) error {
	if h.pingDB == nil {
		return nil
	}
	return h.pingDB(ctx)
}

func (h *serverHealthChecker) PingBlob(ctx context.Context) error {
	if h.pingBlob == nil {
		return nil
	}
	return h.pingBlob(ctx)
}
