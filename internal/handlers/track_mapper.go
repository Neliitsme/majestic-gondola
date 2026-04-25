package handlers

import (
	"majestic-gondola/internal/models"
	"time"
)

func CreateToTrack(tr CreateTrackRequest) (*models.Track, error) {
	track := &models.Track{
		Name:   tr.Name,
		Author: tr.Author,
		Genres: tr.Genres,
	}

	if tr.ReleaseDate != nil && *tr.ReleaseDate != "" {
		releaseDate, err := time.Parse(time.DateTime, *tr.ReleaseDate)

		if err != nil {
			return nil, err
		}

		track.ReleaseDate = releaseDate
	}

	return track, nil
}

func UpdateToTrack(id int, tr UpdateTrackRequest) (*models.Track, error) {
	track := &models.Track{
		Id:     id,
		Name:   tr.Name,
		Author: tr.Author,
		Genres: tr.Genres,
	}

	if tr.ReleaseDate != nil && *tr.ReleaseDate != "" {
		releaseDate, err := time.Parse(time.DateTime, *tr.ReleaseDate)

		if err != nil {
			return nil, err
		}

		track.ReleaseDate = releaseDate
	}

	return track, nil
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
