package handlers

import (
	"majestic-gondola/internal/models"
	"time"
)

func CreateToTrack(tr CreateTrackRequest) (*models.Track, error) {
	releaseDate, err := time.Parse(time.DateTime, tr.ReleaseDate)
	if err != nil {
		return nil, err
	}

	track := &models.Track{
		Name:        tr.Name,
		Author:      tr.Author,
		ReleaseDate: releaseDate,
		Genres:      tr.Genres,
	}

	return track, nil
}

func UpdateToTrack(tr UpdateTrackRequest) (*models.Track, error) {
	releaseDate, err := time.Parse(time.DateTime, tr.ReleaseDate)
	if err != nil {
		return nil, err
	}

	return &models.Track{
		Id:          tr.Id,
		Name:        tr.Name,
		Author:      tr.Author,
		ReleaseDate: releaseDate,
		Genres:      tr.Genres,
	}, nil
}

func ToTrackResponse(t *models.Track) TrackResponse {
	return TrackResponse{
		Id:          t.Id,
		Name:        t.Name,
		Author:      t.Author,
		ReleaseDate: t.ReleaseDate.Format(time.DateTime),
		Genres:      t.Genres,
		CreatedAt:   t.CreatedAt.Format(time.DateTime),
	}
}
