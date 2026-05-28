package inproc

import (
	"context"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/fanboykun/smokery/apps/core/internal/port"
)

// RetentionCleaner periodically deletes old runs and their artifacts.
type RetentionCleaner struct {
	runs      port.RunRepo
	artifacts port.ArtifactRepo
	blob      port.BlobStore
	ttl       time.Duration
	interval  time.Duration
	cancel    context.CancelFunc
}

func NewRetentionCleaner(runs port.RunRepo, artifacts port.ArtifactRepo, blob port.BlobStore, ttl, interval time.Duration) *RetentionCleaner {
	return &RetentionCleaner{runs: runs, artifacts: artifacts, blob: blob, ttl: ttl, interval: interval}
}

func (c *RetentionCleaner) Start() {
	ctx, cancel := context.WithCancel(context.Background())
	c.cancel = cancel
	go c.loop(ctx)
}

func (c *RetentionCleaner) Stop() {
	if c.cancel != nil {
		c.cancel()
	}
}

func (c *RetentionCleaner) loop(ctx context.Context) {
	ticker := time.NewTicker(c.interval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			c.cleanup(ctx)
		}
	}
}

func (c *RetentionCleaner) cleanup(ctx context.Context) {
	before := time.Now().Add(-c.ttl)
	count, err := c.runs.DeleteOlderThan(ctx, before)
	if err != nil {
		log.Error().Err(err).Msg("retention cleanup: failed to delete old runs")
		return
	}
	if count > 0 {
		log.Info().Int("deleted", count).Dur("ttl", c.ttl).Msg("retention cleanup: removed old runs")
	}
}
