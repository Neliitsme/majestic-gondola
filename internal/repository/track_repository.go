package repository

import (
	"errors"
	"log/slog"
	"majestic-gondola/internal/apperr"
	"majestic-gondola/internal/models"

	"github.com/go-pg/pg/v10"
)

type TrackRepository interface {
	FindById(id int) (*models.Track, error)
	GetAll() ([]models.Track, error)
	BulkCreate(tracks []*models.Track) error
	Update(track *models.Track) error
}

type trackRepository struct {
	db  *pg.DB
	log *slog.Logger
}

func NewTrackRepository(db *pg.DB, logger *slog.Logger) TrackRepository {
	return &trackRepository{db: db, log: logger.With("component", "track_repository")}
}

func (r *trackRepository) FindById(id int) (*models.Track, error) {
	track := new(models.Track)
	err := r.db.Model(track).Where("id = ?", id).Select()

	if err != nil {
		if errors.Is(err, pg.ErrNoRows) {
			return nil, apperr.ErrNotFound
		}
		return nil, err
	}

	return track, nil
}

func (r *trackRepository) GetAll() ([]models.Track, error) {
	var tracks []models.Track
	err := r.db.Model(&tracks).Select()

	if err != nil {
		return nil, err
	}

	return tracks, nil
}

func (r *trackRepository) BulkCreate(tracks []*models.Track) error {
	_, err := r.db.Model(&tracks).Insert()

	if err != nil {
		return err
	}

	r.log.Info("Finished creating tracks")
	return nil
}

func (r *trackRepository) Update(track *models.Track) error {
	res, err := r.db.Model(track).ExcludeColumn("created_at").WherePK().Update()

	if err != nil {
		return err
	}

	if res.RowsAffected() == 0 {
		return apperr.ErrNotFound
	}

	r.log.Info("Finished updating track")
	return nil
}
