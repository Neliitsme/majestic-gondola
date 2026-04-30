package repository

import (
	"errors"
	"log/slog"
	"majestic-gondola/internal/apperr"
	"majestic-gondola/internal/models"

	"github.com/go-pg/pg/v10"
)

type ArtistRepository interface {
	FindById(id int) (*models.Artist, error)
	GetAll() ([]models.Artist, error)
	BulkCreate(artists []*models.Artist) error
	Update(artist *models.Artist) error
}

type artistRepository struct {
	db  *pg.DB
	log *slog.Logger
}

func NewArtistRepository(db *pg.DB, logger *slog.Logger) ArtistRepository {
	return &artistRepository{db: db, log: logger.With("component", "artist_repository")}
}

// BulkCreate implements [ArtistRepository].
func (a *artistRepository) BulkCreate(artists []*models.Artist) error {
	_, err := a.db.Model(&artists).Insert()

	if err != nil {
		return err
	}

	a.log.Info("Finished creating artists")
	return nil
}

// FindById implements [ArtistRepository].
func (a *artistRepository) FindById(id int) (*models.Artist, error) {
	artist := new(models.Artist)
	err := a.db.Model(artist).Where("artist_id = ?", id).Select()

	if err != nil {
		if errors.Is(err, pg.ErrNoRows) {
			return nil, apperr.ErrNotFound
		}
		return nil, err
	}

	return artist, nil
}

// GetAll implements [ArtistRepository].
func (a *artistRepository) GetAll() ([]models.Artist, error) {
	var artists []models.Artist
	err := a.db.Model(&artists).Select()

	if err != nil {
		return nil, err
	}

	return artists, nil
}

// Update implements [ArtistRepository].
func (a *artistRepository) Update(artist *models.Artist) error {
	res, err := a.db.Model(artist).ExcludeColumn("created_at").WherePK().Update()

	if err != nil {
		return err
	}

	if res.RowsAffected() == 0 {
		return apperr.ErrNotFound
	}

	a.log.Info("Finished updating artist")
	return nil
}
