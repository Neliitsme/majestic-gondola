package service

import (
	"errors"
	"fmt"
	"log/slog"
	"majestic-gondola/internal/apperr"
	"majestic-gondola/internal/models"
	"majestic-gondola/internal/repository"
	"time"
)

type TrackService interface {
	Get(id int) (*models.Track, error)
	GetAll() ([]models.Track, error)
	BulkCreate(tracks []*models.Track) error
	Update(track *models.Track) error
	Generate(count int, artistId *int) error
}

type trackService struct {
	trackRepository repository.TrackRepository
	log             *slog.Logger
}

func NewTrackService(trackRepository repository.TrackRepository, logger *slog.Logger) TrackService {
	return &trackService{log: logger.With("component", "track_service"), trackRepository: trackRepository}
}

func (s *trackService) Get(id int) (*models.Track, error) {
	track, err := s.trackRepository.FindById(id)

	if err != nil {
		if errors.Is(err, apperr.ErrNotFound) {
			return nil, apperr.NotFound("Track not found", err)
		}
		s.log.Error("Failed to fetch a track by id", slog.Any("error", err), slog.Int("id", id))
		return nil, apperr.Internal(err)
	}

	return track, nil
}

func (s *trackService) GetAll() ([]models.Track, error) {
	tracks, err := s.trackRepository.GetAll()

	if err != nil {
		s.log.Error("Failed to fetch tracks", slog.Any("error", err))
		return nil, apperr.Internal(err)
	}

	return tracks, nil
}

func (s *trackService) BulkCreate(tracks []*models.Track) error {
	err := s.trackRepository.BulkCreate(tracks)

	if err != nil {
		s.log.Error("Failed to bulk create tracks", slog.Any("error", err), slog.Int("parsed_tracks", len(tracks)))
		return apperr.Internal(err)
	}

	s.log.Info("Bulk created tracks", slog.Int("count", len(tracks)))
	return nil
}

func (s *trackService) Update(track *models.Track) error {
	err := s.trackRepository.Update(track)

	if err != nil {
		if errors.Is(err, apperr.ErrNotFound) {
			return apperr.NotFound("Track not found", err)
		}
		s.log.Error("Failed to update the track", slog.Any("error", err), slog.Any("parsed_track", track))
		return apperr.Internal(err)
	}

	s.log.Info("Updated a track")
	return nil
}

func (s *trackService) Generate(count int, artistId *int) error {

	tracks := make([]*models.Track, 0, count)
	for i := range count {
		track := &models.Track{
			Name:        fmt.Sprintf("Track %d", i),
			ArtistId:    artistId,
			ReleaseDate: time.Now(),
			Genres:      []string{"Tag"},
		}

		tracks = append(tracks, track)
	}

	return s.BulkCreate(tracks)
}
