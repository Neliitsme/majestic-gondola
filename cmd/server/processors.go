package main

import (
	"context"
	"log/slog"

	"github.com/robfig/cron/v3"

	"majestic-gondola/bootstrap"
	"majestic-gondola/internal/processor"
	"majestic-gondola/internal/repository"
)

func setupProcessors(
	c *cron.Cron,
	ctx context.Context,
	cfg *bootstrap.Config,
	rr repository.ReviewRepository,
	tr repository.TrackRepository,
	ar repository.ArtistRepository,
	sc repository.ScoreCommitter,
	logger *slog.Logger,
) {
	reviewProc := processor.NewReviewProcessor(rr, sc, logger)
	addJob(c, cfg.ReviewProcessorSchedule, "review", logger, func() {
		if err := reviewProc.Run(ctx); err != nil {
			logger.Error("Review processor failed", slog.Any("error", err))
		}
	})

	artistProc := processor.NewArtistProcessor(tr, ar, logger)
	addJob(c, cfg.ArtistProcessorSchedule, "artist", logger, func() {
		if err := artistProc.Run(ctx); err != nil {
			logger.Error("Artist processor failed", slog.Any("error", err))
		}
	})
}

func addJob(c *cron.Cron, schedule, name string, logger *slog.Logger, fn func()) {
	if schedule == "" {
		return
	}
	if _, err := c.AddFunc(schedule, fn); err != nil {
		logger.Error("Invalid schedule", slog.String("job", name), slog.String("schedule", schedule), slog.Any("error", err))
	} else {
		logger.Info("Registered job", slog.String("job", name), slog.String("schedule", schedule))
	}
}
