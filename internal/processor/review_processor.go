package processor

import (
	"context"
	"log/slog"
	"majestic-gondola/internal/repository"
	"sync/atomic"
)

type ReviewProcessor struct {
	reviewRepo repository.ReviewRepository
	trackRepo  repository.TrackRepository
	artistRepo repository.ArtistRepository
	committer  repository.ScoreCommitter
	log        *slog.Logger
	isRunning  atomic.Bool
}

func NewReviewProcessor(
	reviewRepo repository.ReviewRepository,
	trackRepo repository.TrackRepository,
	artistRepo repository.ArtistRepository,
	committer repository.ScoreCommitter,
	log *slog.Logger,
) *ReviewProcessor {
	return &ReviewProcessor{
		reviewRepo: reviewRepo,
		trackRepo:  trackRepo,
		artistRepo: artistRepo,
		committer:  committer,
		log:        log.With("component", "review_processor"),
	}
}

type trackData struct {
	sum          int
	count        int
	currentScore int
}

func (p *ReviewProcessor) Run(ctx context.Context) error {
	if !p.isRunning.CompareAndSwap(false, true) {
		p.log.Warn("processor already running, skipping")
		return nil
	}
	defer p.isRunning.Store(false)

	reviews, err := p.reviewRepo.GetUnprocessed()
	if err != nil {
		return err
	}
	if len(reviews) == 0 {
		return nil
	}

	data := make(map[int]trackData)
	var ids []int
	for _, r := range reviews {
		if r.TrackId == nil || r.Track == nil {
			continue
		}
		d := data[*r.TrackId]
		d.sum += r.Score
		d.count++
		d.currentScore = r.Track.Score
		data[*r.TrackId] = d
		ids = append(ids, r.Id)
	}

	if len(data) == 0 {
		return nil
	}

	// TODO: Store review count in the track to calculate correctly
	newScores := make(map[int]int, len(data))
	for trackId, d := range data {
		newScores[trackId] = (d.currentScore + d.sum/d.count) / 2
	}

	if err := p.committer.CommitBatch(ctx, newScores, ids); err != nil {
		return err
	}

	p.log.Info("processed reviews", slog.Int("reviews", len(ids)), slog.Int("tracks", len(newScores)))
	return nil
}
