package service

import (
	"log/slog"
	"majestic-gondola/internal/models"
	"majestic-gondola/internal/repository"
)

type TrackService interface {
	Get(id int) (*models.Track, error)
	GetAll() ([]models.Track, error)
	BulkCreate(tracks []*models.Track) error
	Update(track *models.Track) error
}

type trackService struct {
	trackRepository repository.TrackRepository
	log             *slog.Logger
}

func NewTrackService(trackRepository repository.TrackRepository, logger *slog.Logger) TrackService {
	return &trackService{log: logger.With("component", "track_service"), trackRepository: trackRepository}
}

func (s *trackService) Get(id int) (*models.Track, error) {
	return s.trackRepository.FindById(id)
}
func (s *trackService) GetAll() ([]models.Track, error) {
	return s.trackRepository.GetAll()
}
func (s *trackService) BulkCreate(tracks []*models.Track) error {
	return s.trackRepository.BulkCreate(tracks)
}
func (s *trackService) Update(track *models.Track) error {
	return s.trackRepository.Update(track)
}
