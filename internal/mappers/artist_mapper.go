package mappers

import (
	"majestic-gondola/internal/dto"
	"majestic-gondola/internal/models"
	"time"
)

func ToArtistResponse(a *models.Artist) dto.ArtistResponse {
	return dto.ArtistResponse{
		Id:        a.Id,
		Name:      a.Name,
		Score:     a.Score,
		CreatedAt: a.CreatedAt.Format(time.DateTime),
	}
}

func ToArtistResponseList(arts []models.Artist) []dto.ArtistResponse {
	responses := make([]dto.ArtistResponse, 0, len(arts))
	for i := range arts {
		responses = append(responses, ToArtistResponse(&arts[i]))
	}
	return responses
}

func CreateToArtist(ar dto.CreateArtistRequest) *models.Artist {
	return &models.Artist{
		Name: ar.Name,
	}
}

func CreateToArtistList(ars []dto.CreateArtistRequest) []*models.Artist {
	artists := make([]*models.Artist, 0, len(ars))
	for i := range ars {
		artist := CreateToArtist(ars[i])
		artists = append(artists, artist)
	}
	return artists
}

func UpdateToArtist(id int, ar dto.UpdateArtistRequest) *models.Artist {
	return &models.Artist{
		Id:   id,
		Name: ar.Name,
	}
}
