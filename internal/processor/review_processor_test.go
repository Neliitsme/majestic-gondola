package processor_test

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"testing"

	"majestic-gondola/internal/models"
	"majestic-gondola/internal/processor"
	"majestic-gondola/internal/repository/mocks"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func discardLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(io.Discard, nil))
}

func newProc(t *testing.T) (
	*processor.ReviewProcessor,
	*mocks.MockReviewRepository,
	*mocks.MockTrackRepository,
	*mocks.MockArtistRepository,
	*mocks.MockScoreCommitter,
) {
	rr := mocks.NewMockReviewRepository(t)
	tr := mocks.NewMockTrackRepository(t)
	ar := mocks.NewMockArtistRepository(t)
	committer := mocks.NewMockScoreCommitter(t)
	proc := processor.NewReviewProcessor(rr, tr, ar, committer, discardLogger())
	return proc, rr, tr, ar, committer
}

func TestRun_NoReviews(t *testing.T) {
	proc, rr, _, _, _ := newProc(t)
	rr.EXPECT().GetUnprocessed().Return([]models.Review{}, nil)

	require.NoError(t, proc.Run(context.Background()))
}

func TestRun_AllNilTrackId(t *testing.T) {
	proc, rr, _, _, _ := newProc(t)
	rr.EXPECT().GetUnprocessed().Return([]models.Review{
		{Id: 1, TrackId: nil, Score: 80},
		{Id: 2, TrackId: nil, Score: 60},
	}, nil)

	require.NoError(t, proc.Run(context.Background()))
}

func TestRun_AllNilTrack(t *testing.T) {
	proc, rr, _, _, _ := newProc(t)
	rr.EXPECT().GetUnprocessed().Return([]models.Review{
		{Id: 1, TrackId: new(10), Track: nil, Score: 80},
	}, nil)

	require.NoError(t, proc.Run(context.Background()))
}

func TestRun_SingleReview(t *testing.T) {
	// track score=50, review score=70 => avgNew=70 => newScore=(50+70)/2=60
	proc, rr, _, _, committer := newProc(t)

	rr.EXPECT().GetUnprocessed().Return([]models.Review{
		{Id: 1, TrackId: new(10), Score: 70, Track: &models.Track{Id: 10, Score: 50}},
	}, nil)
	committer.EXPECT().CommitBatch(mock.Anything, map[int]int{10: 60}, []int{1}).Return(nil)

	require.NoError(t, proc.Run(context.Background()))
}

func TestRun_MultipleReviewsSameTrack(t *testing.T) {
	proc, rr, _, _, committer := newProc(t)

	track := &models.Track{Id: 10, Score: 50}
	rr.EXPECT().GetUnprocessed().Return([]models.Review{
		{Id: 1, TrackId: new(10), Score: 60, Track: track},
		{Id: 2, TrackId: new(10), Score: 80, Track: track},
	}, nil)
	committer.EXPECT().CommitBatch(mock.Anything, map[int]int{10: 60}, []int{1, 2}).Return(nil)

	require.NoError(t, proc.Run(context.Background()))
}

func TestRun_MultipleTrackIds(t *testing.T) {
	proc, rr, _, _, committer := newProc(t)

	rr.EXPECT().GetUnprocessed().Return([]models.Review{
		{Id: 1, TrackId: new(1), Score: 70, Track: &models.Track{Id: 1, Score: 50}},
		{Id: 2, TrackId: new(2), Score: 40, Track: &models.Track{Id: 2, Score: 20}},
	}, nil)

	var capturedScores map[int]int
	var capturedIds []int
	committer.EXPECT().CommitBatch(mock.Anything, mock.Anything, mock.Anything).
		RunAndReturn(func(_ context.Context, scores map[int]int, ids []int) error {
			capturedScores = scores
			capturedIds = ids
			return nil
		})

	require.NoError(t, proc.Run(context.Background()))
	assert.Equal(t, map[int]int{1: 60, 2: 30}, capturedScores)
	assert.ElementsMatch(t, []int{1, 2}, capturedIds)
}

func TestRun_GetUnprocessedError(t *testing.T) {
	proc, rr, _, _, _ := newProc(t)
	wantErr := errors.New("db error")
	rr.EXPECT().GetUnprocessed().Return(nil, wantErr)

	assert.ErrorIs(t, proc.Run(context.Background()), wantErr)
}

func TestRun_CommitBatchError(t *testing.T) {
	proc, rr, _, _, committer := newProc(t)
	wantErr := errors.New("tx error")

	rr.EXPECT().GetUnprocessed().Return([]models.Review{
		{Id: 1, TrackId: new(1), Score: 70, Track: &models.Track{Id: 1, Score: 50}},
	}, nil)
	committer.EXPECT().CommitBatch(mock.Anything, mock.Anything, mock.Anything).Return(wantErr)

	assert.ErrorIs(t, proc.Run(context.Background()), wantErr)
}
