package repository

import (
	"context"
	"majestic-gondola/internal/models"

	"github.com/go-pg/pg/v10"
)

type ScoreCommitter interface {
	CommitBatch(ctx context.Context, trackScores map[int]TrackScoresUpdate, reviewIds []int) error
}

type scoreCommitter struct {
	db *pg.DB
}

func NewScoreCommitter(db *pg.DB) ScoreCommitter {
	return &scoreCommitter{db: db}
}

type TrackScoresUpdate struct {
	Score int
	Count int
}

func (c *scoreCommitter) CommitBatch(ctx context.Context, trackData map[int]TrackScoresUpdate, reviewIds []int) error {
	return c.db.RunInTransaction(ctx, func(tx *pg.Tx) error {
		for id, data := range trackData {
			if _, err := tx.Model(new(models.Track)).
				Set("score = ?", data.Score).
				Set("review_count = ?", data.Count).
				Where("track_id = ?", id).
				Update(); err != nil {
				return err
			}
		}
		_, err := tx.Model(new(models.Review)).
			Set("is_processed = true").
			Where("review_id IN (?)", pg.In(reviewIds)).
			Update()
		return err
	})
}
