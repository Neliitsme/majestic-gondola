package service

import (
	"fmt"
	"log/slog"
	"majestic-gondola/internal/models"
	"majestic-gondola/internal/repository"
	"time"
)

type TrackService interface {
	Get(id int) (*models.Track, error)
	GetAll() ([]models.Track, error)
	BulkCreate(tracks []*models.Track) error
	Update(track *models.Track) error
	Generate(count int) error
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
		s.log.Error("Failed to fetch a track by id", slog.Any("error", err), slog.Int("id", id))
	}

	return track, err
}

func (s *trackService) GetAll() ([]models.Track, error) {
	return s.trackRepository.GetAll()
}

func (s *trackService) BulkCreate(tracks []*models.Track) error {
	err := s.trackRepository.BulkCreate(tracks)

	if err != nil {
		s.log.Error("Failed to bulk create tracks", slog.Any("error", err), slog.Int("parsed_tracks", len(tracks)))
	} else {
		s.log.Info("Created a new track")
	}

	return err
}

func (s *trackService) Update(track *models.Track) error {
	err := s.trackRepository.Update(track)

	if err != nil {
		s.log.Error("Failed to update the track", slog.Any("error", err), slog.Any("parsed_track", track))
	} else {
		s.log.Info("Updated a track")
	}

	return err
}

func (s *trackService) Generate(count int) error {

	tracks := make([]*models.Track, 0, count)
	for i := range count {
		track := &models.Track{
			Name:        fmt.Sprintf("Track %d", i),
			Author:      fmt.Sprintf("Author %d", i),
			ReleaseDate: time.Now(),
			Genres:      []string{"Tag"},
		}

		tracks = append(tracks, track)
	}

	err := s.trackRepository.BulkCreate(tracks)

	if err != nil {
		s.log.Error("Failed to bulk create tracks during population", slog.Any("error", err), slog.Int("generated_tracks", len(tracks)))
	} else {
		s.log.Info("Created several new tracks")
	}

	return err
}
