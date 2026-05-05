package mappers

import (
	"majestic-gondola/internal/dto"
	"majestic-gondola/internal/models"
	"time"
)

func CreateToTrack(tr dto.CreateTrackRequest) (*models.Track, error) {
	track := &models.Track{
		Name:     tr.Name,
		ArtistId: tr.ArtistId,
		Genres:   tr.Genres,
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

func UpdateToTrack(id int, tr dto.UpdateTrackRequest) (*models.Track, error) {
	track := &models.Track{
		Id:       id,
		Name:     tr.Name,
		ArtistId: tr.ArtistId,
		Genres:   tr.Genres,
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

func ToTrackResponse(t *models.Track) dto.TrackResponse {
	return dto.TrackResponse{
		Id:          t.Id,
		Name:        t.Name,
		ArtistId:    t.ArtistId,
		ReleaseDate: t.ReleaseDate.Format(time.DateTime),
		Genres:      t.Genres,
		CreatedAt:   t.CreatedAt.Format(time.DateTime),
		Score:       t.Score,
		ReviewCount: t.ReviewCount,
	}
}

func ToTrackResponseList(ts []models.Track) []dto.TrackResponse {
	responses := make([]dto.TrackResponse, 0, len(ts))
	for i := range ts {
		responses = append(responses, ToTrackResponse(&ts[i]))
	}
	return responses
}
