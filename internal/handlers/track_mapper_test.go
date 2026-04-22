package handlers_test

import (
	"majestic-gondola/internal/handlers"
	"majestic-gondola/internal/models"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestToTrackResponse(t *testing.T) {
	assert := assert.New(t)

	tests := []struct {
		name                string
		input               models.Track
		expectedReleaseDate string
	}{
		{
			name: "Valid full model",
			input: models.Track{
				Id:          1,
				Name:        "Test",
				Author:      "Test",
				ReleaseDate: time.Date(2006, 01, 02, 15, 04, 05, 06, time.UTC),
				Genres:      []string{"Test"},
				CreatedAt:   time.Now(),
			},
			expectedReleaseDate: "2006-01-02 15:04:05",
		},
		{
			name: "Valid min model",
			input: models.Track{
				Id:   1,
				Name: "Test",
			},
			expectedReleaseDate: "0001-01-01 00:00:00",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			track := handlers.ToTrackResponse(&test.input)
			require.NotNil(t, track)
			assert.Equal(test.input.Id, track.Id)
			assert.Equal(test.input.Name, track.Name)
			assert.Equal(test.expectedReleaseDate, track.ReleaseDate)
		})
	}
}

func TestCreateToTrack(t *testing.T) {
	assert := assert.New(t)

	tests := []struct {
		name    string
		input   handlers.CreateTrackRequest
		wantErr bool
	}{
		{
			name: "Valid request",
			input: handlers.CreateTrackRequest{
				Name:        "Song",
				ReleaseDate: new("2026-04-20 20:00:00"),
			},
			wantErr: false,
		},
		{
			name: "Invalid date format",
			input: handlers.CreateTrackRequest{
				Name:        "Song",
				ReleaseDate: new("invalid-date"),
			},
			wantErr: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			track, err := handlers.CreateToTrack(test.input)
			if test.wantErr {
				assert.Error(err)
				assert.Nil(track)
			} else {
				assert.NoError(err)
				assert.NotNil(track)
			}
		})
	}
}

func TestUpdateToTrack(t *testing.T) {
	assert := assert.New(t)

	tests := []struct {
		name    string
		input   handlers.UpdateTrackRequest
		wantErr bool
	}{
		{
			name: "Valid request",
			input: handlers.UpdateTrackRequest{
				Name:        "Song",
				ReleaseDate: new("2026-04-20 20:00:00"),
			},
			wantErr: false,
		},
		{
			name: "Invalid date format",
			input: handlers.UpdateTrackRequest{
				Name:        "Song",
				ReleaseDate: new("invalid-date"),
			},
			wantErr: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			track, err := handlers.UpdateToTrack(test.input)
			if test.wantErr {
				assert.Error(err)
				assert.Nil(track)
			} else {
				assert.NoError(err)
				assert.NotNil(track)
			}
		})
	}
}
