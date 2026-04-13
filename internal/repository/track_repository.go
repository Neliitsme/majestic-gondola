package repository

import (
	"log/slog"
	"majestic-gondola/internal/models"

	"github.com/go-pg/pg/v10"
)

type TrackStore interface {
	FindById(id int) (*models.Track, error)
	GetAll() ([]models.Track, error)
	BulkCreate(tracks []*models.Track) error
	Update(track *models.Track) error
}

type TrackRepository struct {
	db  *pg.DB
	log *slog.Logger
}

func NewTrackRepository(db *pg.DB, logger *slog.Logger) *TrackRepository {
	return &TrackRepository{db: db, log: logger.With("component", "track_repository")}
}

func (r *TrackRepository) FindById(id int) (*models.Track, error) {
	track := new(models.Track)
	err := r.db.Model(track).Where("id = ?", id).Select()
	return track, err
}

func (r *TrackRepository) GetAll() ([]models.Track, error) {
	var tracks []models.Track
	err := r.db.Model(&tracks).Select()
	return tracks, err
}

func (r *TrackRepository) BulkCreate(tracks []*models.Track) error {
	_, err := r.db.Model(&tracks).Insert()
	r.log.Info("Finished BulkCreate")
	return err
}

func (r *TrackRepository) Update(track *models.Track) error {
	_, err := r.db.Model(track).ExcludeColumn("created_at").WherePK().Update()
	r.log.Info("Finished Update")
	return err
}
