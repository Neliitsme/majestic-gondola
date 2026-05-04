package repository

import (
	"context"
	"majestic-gondola/internal/models"

	"github.com/go-pg/pg/v10"
)

type ScoreCommitter interface {
	CommitBatch(ctx context.Context, trackScores map[int]int, reviewIds []int) error
}

type scoreCommitter struct {
	db *pg.DB
}

func NewScoreCommitter(db *pg.DB) ScoreCommitter {
	return &scoreCommitter{db: db}
}

func (c *scoreCommitter) CommitBatch(ctx context.Context, trackScores map[int]int, reviewIds []int) error {
	return c.db.RunInTransaction(ctx, func(tx *pg.Tx) error {
		for id, score := range trackScores {
			if _, err := tx.Model(new(models.Track)).
				Set("score = ?", score).
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
