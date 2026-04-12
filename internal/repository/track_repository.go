package repository

import (
	"log/slog"
	"majestic-gondola/internal/models"

	"github.com/go-pg/pg/v10"
)

type TrackRepository struct {
	db  *pg.DB
	log *slog.Logger
}

func NewTrackRepository(db *pg.DB, logger *slog.Logger) *TrackRepository {
	return &TrackRepository{db: db, log: logger}
}

func (r *TrackRepository) FindById(id int) (*models.Track, error) {
	track := new(models.Track)
	err := r.db.Model(track).Where("id = ?", id).Select()
	r.log.Info("Finished FindById")
	return track, err
}

func (r *TrackRepository) GetAll() ([]models.Track, error) {
	var tracks []models.Track
	err := r.db.Model(&tracks).Select()
	r.log.Info("Finished GetAll")
	return tracks, err
}

func (r *TrackRepository) BulkCreate(tracks []*models.Track) error {
	_, err := r.db.Model(&tracks).Insert()
	r.log.Info("Finished BulkCreate")
	return err
}
