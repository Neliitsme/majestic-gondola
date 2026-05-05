package processor

import (
	"context"
	"log/slog"
	"majestic-gondola/internal/repository"
	"sync/atomic"
)

type ArtistProcessor struct {
	trackRepo  repository.TrackRepository
	artistRepo repository.ArtistRepository
	log        *slog.Logger
	isRunning  atomic.Bool
}

func NewArtistProcessor(
	trackRepo repository.TrackRepository,
	artistRepo repository.ArtistRepository,
	log *slog.Logger,
) *ArtistProcessor {
	return &ArtistProcessor{
		trackRepo:  trackRepo,
		artistRepo: artistRepo,
		log:        log.With("component", "artist_processor"),
	}
}

type artistData struct {
	total int
	count int
}

func (p *ArtistProcessor) Run(ctx context.Context) error {
	if !p.isRunning.CompareAndSwap(false, true) {
		p.log.Warn("Artist processor already running, skipping")
		return nil
	}
	defer p.isRunning.Store(false)

	tracks, err := p.trackRepo.GetAll()
	if err != nil {
		return err
	}

	data := make(map[int]artistData)
	for _, t := range tracks {
		if t.ArtistId == nil || t.ReviewCount == 0 {
			continue
		}
		d := data[*t.ArtistId]
		d.total += t.Score
		d.count++
		data[*t.ArtistId] = d
	}

	if len(data) == 0 {
		return nil
	}

	scores := make(map[int]int, len(data))
	for artistId, d := range data {
		scores[artistId] = d.total / d.count
	}

	if err := p.artistRepo.BulkUpdateScores(ctx, scores); err != nil {
		return err
	}

	p.log.Info("Updated artist scores", slog.Int("artists", len(scores)))
	return nil
}
