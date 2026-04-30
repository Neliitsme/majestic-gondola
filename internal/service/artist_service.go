package service

import (
	"errors"
	"log/slog"
	"majestic-gondola/internal/apperr"
	"majestic-gondola/internal/models"
	"majestic-gondola/internal/repository"
)

type ArtistService interface {
	Get(id int) (*models.Artist, error)
	GetAll() ([]models.Artist, error)
	BulkCreate(artists []*models.Artist) error
	Update(artist *models.Artist) error
}

type artistService struct {
	artistRepository repository.ArtistRepository
	log              *slog.Logger
}

func NewArtistService(artistRepository repository.ArtistRepository, logger *slog.Logger) ArtistService {
	return &artistService{log: logger.With("component", "artist_service"), artistRepository: artistRepository}
}

// BulkCreate implements [ArtistService].
func (a *artistService) BulkCreate(artists []*models.Artist) error {
	err := a.artistRepository.BulkCreate(artists)

	if err != nil {
		a.log.Error("Failed to bulk create artists", slog.Any("error", err), slog.Int("parsed_artists", len(artists)))
		return apperr.Internal(err)
	}

	a.log.Info("Bulk created artists", slog.Int("count", len(artists)))
	return nil
}

// Get implements [ArtistService].
func (a *artistService) Get(id int) (*models.Artist, error) {
	artist, err := a.artistRepository.FindById(id)

	if err != nil {
		if errors.Is(err, apperr.ErrNotFound) {
			return nil, apperr.NotFound("Artist not found", err)
		}
		a.log.Error("Failed to fetch a artist by id", slog.Any("error", err), slog.Int("id", id))
		return nil, apperr.Internal(err)
	}

	return artist, nil
}

// GetAll implements [ArtistService].
func (a *artistService) GetAll() ([]models.Artist, error) {
	artists, err := a.artistRepository.GetAll()

	if err != nil {
		a.log.Error("Failed to fetch artists", slog.Any("error", err))
		return nil, apperr.Internal(err)
	}

	return artists, nil
}

// Update implements [ArtistService].
func (a *artistService) Update(artist *models.Artist) error {
	err := a.artistRepository.Update(artist)

	if err != nil {
		if errors.Is(err, apperr.ErrNotFound) {
			return apperr.NotFound("Artist not found", err)
		}
		a.log.Error("Failed to update the artist", slog.Any("error", err), slog.Any("parsed_artist", artist))
		return apperr.Internal(err)
	}

	a.log.Info("Updated a artist")
	return nil
}
