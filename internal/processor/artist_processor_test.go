package processor_test

import (
	"context"
	"errors"
	"testing"

	"majestic-gondola/internal/models"
	"majestic-gondola/internal/processor"
	"majestic-gondola/internal/repository/mocks"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func newArtistProc(t *testing.T) (
	*processor.ArtistProcessor,
	*mocks.MockTrackRepository,
	*mocks.MockArtistRepository,
) {
	tr := mocks.NewMockTrackRepository(t)
	ar := mocks.NewMockArtistRepository(t)
	proc := processor.NewArtistProcessor(tr, ar, discardLogger())
	return proc, tr, ar
}

func TestArtistRun_NoTracks(t *testing.T) {
	proc, tr, _ := newArtistProc(t)
	tr.EXPECT().GetAll().Return([]models.Track{}, nil)

	require.NoError(t, proc.Run(context.Background()))
}

func TestArtistRun_AllNilArtistId(t *testing.T) {
	proc, tr, _ := newArtistProc(t)
	tr.EXPECT().GetAll().Return([]models.Track{
		{Id: 1, ArtistId: nil, Score: 80, ReviewCount: 2},
	}, nil)

	require.NoError(t, proc.Run(context.Background()))
}

func TestArtistRun_AllZeroReviewCount(t *testing.T) {
	proc, tr, _ := newArtistProc(t)
	tr.EXPECT().GetAll().Return([]models.Track{
		{Id: 1, ArtistId: new(1), Score: 80, ReviewCount: 0},
	}, nil)

	require.NoError(t, proc.Run(context.Background()))
}

func TestArtistRun_SingleTrack(t *testing.T) {
	proc, tr, ar := newArtistProc(t)
	tr.EXPECT().GetAll().Return([]models.Track{
		{Id: 1, ArtistId: new(10), Score: 70, ReviewCount: 2},
	}, nil)
	ar.EXPECT().BulkUpdateScores(mock.Anything, map[int]int{10: 70}).Return(nil)

	require.NoError(t, proc.Run(context.Background()))
}

func TestArtistRun_MultipleTracksSameArtist(t *testing.T) {
	proc, tr, ar := newArtistProc(t)
	tr.EXPECT().GetAll().Return([]models.Track{
		{Id: 1, ArtistId: new(10), Score: 60, ReviewCount: 1},
		{Id: 2, ArtistId: new(10), Score: 80, ReviewCount: 3},
	}, nil)
	ar.EXPECT().BulkUpdateScores(mock.Anything, map[int]int{10: 70}).Return(nil)

	require.NoError(t, proc.Run(context.Background()))
}

func TestArtistRun_ZeroReviewTrackExcluded(t *testing.T) {
	proc, tr, ar := newArtistProc(t)
	tr.EXPECT().GetAll().Return([]models.Track{
		{Id: 1, ArtistId: new(10), Score: 60, ReviewCount: 2},
		{Id: 2, ArtistId: new(10), Score: 100, ReviewCount: 0},
	}, nil)
	ar.EXPECT().BulkUpdateScores(mock.Anything, map[int]int{10: 60}).Return(nil)

	require.NoError(t, proc.Run(context.Background()))
}

func TestArtistRun_MultipleArtists(t *testing.T) {
	proc, tr, ar := newArtistProc(t)
	tr.EXPECT().GetAll().Return([]models.Track{
		{Id: 1, ArtistId: new(1), Score: 60, ReviewCount: 1},
		{Id: 2, ArtistId: new(2), Score: 80, ReviewCount: 2},
	}, nil)

	var captured map[int]int
	ar.EXPECT().BulkUpdateScores(mock.Anything, mock.Anything).
		RunAndReturn(func(_ context.Context, scores map[int]int) error {
			captured = scores
			return nil
		})

	require.NoError(t, proc.Run(context.Background()))
	assert.Equal(t, map[int]int{1: 60, 2: 80}, captured)
}

func TestArtistRun_GetAllError(t *testing.T) {
	proc, tr, _ := newArtistProc(t)
	wantErr := errors.New("db error")
	tr.EXPECT().GetAll().Return(nil, wantErr)

	assert.ErrorIs(t, proc.Run(context.Background()), wantErr)
}

func TestArtistRun_BulkUpdateScoresError(t *testing.T) {
	proc, tr, ar := newArtistProc(t)
	wantErr := errors.New("tx error")
	tr.EXPECT().GetAll().Return([]models.Track{
		{Id: 1, ArtistId: new(10), Score: 70, ReviewCount: 1},
	}, nil)
	ar.EXPECT().BulkUpdateScores(mock.Anything, mock.Anything).Return(wantErr)

	assert.ErrorIs(t, proc.Run(context.Background()), wantErr)
}

func TestArtistRun_AlreadyRunning(t *testing.T) {
	proc, tr, _ := newArtistProc(t)

	entered := make(chan struct{})
	release := make(chan struct{})

	tr.EXPECT().GetAll().RunAndReturn(func() ([]models.Track, error) {
		close(entered)
		<-release
		return nil, nil
	})

	done := make(chan error, 1)
	go func() { done <- proc.Run(context.Background()) }()

	<-entered

	err := proc.Run(context.Background())
	require.NoError(t, err)

	close(release)
	require.NoError(t, <-done)
}
