package processor_test

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"testing"

	"majestic-gondola/internal/models"
	"majestic-gondola/internal/processor"
	"majestic-gondola/internal/repository"
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
	*mocks.MockScoreCommitter,
) {
	rr := mocks.NewMockReviewRepository(t)
	tr := mocks.NewMockTrackRepository(t)
	committer := mocks.NewMockScoreCommitter(t)
	proc := processor.NewReviewProcessor(rr, tr, committer, discardLogger())
	return proc, rr, committer
}

func TestReviewRun_NoReviews(t *testing.T) {
	proc, rr, _ := newProc(t)
	rr.EXPECT().GetUnprocessed().Return([]models.Review{}, nil)

	require.NoError(t, proc.Run(context.Background()))
}

func TestReviewRun_AllNilTrackId(t *testing.T) {
	proc, rr, _ := newProc(t)
	rr.EXPECT().GetUnprocessed().Return([]models.Review{
		{Id: 1, TrackId: nil, Score: 80},
		{Id: 2, TrackId: nil, Score: 60},
	}, nil)

	require.NoError(t, proc.Run(context.Background()))
}

func TestReviewRun_AllNilTrack(t *testing.T) {
	proc, rr, _ := newProc(t)
	rr.EXPECT().GetUnprocessed().Return([]models.Review{
		{Id: 1, TrackId: new(10), Track: nil, Score: 80},
	}, nil)

	require.NoError(t, proc.Run(context.Background()))
}

func TestReviewRun_SingleReview(t *testing.T) {
	proc, rr, committer := newProc(t)

	rr.EXPECT().GetUnprocessed().Return([]models.Review{
		{Id: 1, TrackId: new(10), Score: 70, Track: &models.Track{Id: 10, Score: 50, ReviewCount: 1}},
	}, nil)
	committer.EXPECT().CommitBatch(mock.Anything, map[int]repository.TrackScoresUpdate{10: {Score: 60, Count: 2}}, []int{1}).Return(nil)

	require.NoError(t, proc.Run(context.Background()))
}

func TestReviewRun_MultipleReviewsSameTrack(t *testing.T) {
	proc, rr, committer := newProc(t)

	track := &models.Track{Id: 10, Score: 50, ReviewCount: 1}
	rr.EXPECT().GetUnprocessed().Return([]models.Review{
		{Id: 1, TrackId: new(10), Score: 60, Track: track},
		{Id: 2, TrackId: new(10), Score: 70, Track: track},
	}, nil)
	committer.EXPECT().CommitBatch(mock.Anything, map[int]repository.TrackScoresUpdate{10: {Score: 60, Count: 3}}, []int{1, 2}).Return(nil)

	require.NoError(t, proc.Run(context.Background()))
}

func TestReviewRun_MultipleTrackIds(t *testing.T) {
	proc, rr, committer := newProc(t)

	rr.EXPECT().GetUnprocessed().Return([]models.Review{
		{Id: 1, TrackId: new(1), Score: 70, Track: &models.Track{Id: 1, Score: 50, ReviewCount: 1}},
		{Id: 2, TrackId: new(2), Score: 40, Track: &models.Track{Id: 2, Score: 20, ReviewCount: 1}},
	}, nil)

	var capturedScores map[int]repository.TrackScoresUpdate
	var capturedIds []int
	committer.EXPECT().CommitBatch(mock.Anything, mock.Anything, mock.Anything).
		RunAndReturn(func(_ context.Context, scores map[int]repository.TrackScoresUpdate, ids []int) error {
			capturedScores = scores
			capturedIds = ids
			return nil
		})

	require.NoError(t, proc.Run(context.Background()))
	assert.Equal(t, map[int]repository.TrackScoresUpdate{
		1: {Score: 60, Count: 2},
		2: {Score: 30, Count: 2}},
		capturedScores)
	assert.ElementsMatch(t, []int{1, 2}, capturedIds)
}

func TestReviewRun_GetUnprocessedError(t *testing.T) {
	proc, rr, _ := newProc(t)
	wantErr := errors.New("db error")
	rr.EXPECT().GetUnprocessed().Return(nil, wantErr)

	assert.ErrorIs(t, proc.Run(context.Background()), wantErr)
}

func TestReviewRun_CommitBatchError(t *testing.T) {
	proc, rr, committer := newProc(t)
	wantErr := errors.New("tx error")

	rr.EXPECT().GetUnprocessed().Return([]models.Review{
		{Id: 1, TrackId: new(1), Score: 70, Track: &models.Track{Id: 1, Score: 50}},
	}, nil)
	committer.EXPECT().CommitBatch(mock.Anything, mock.Anything, mock.Anything).Return(wantErr)

	assert.ErrorIs(t, proc.Run(context.Background()), wantErr)
}

func TestReviewRun_AlreadyRunning(t *testing.T) {
	proc, rr, _ := newProc(t)

	entered := make(chan struct{})
	release := make(chan struct{})

	rr.EXPECT().GetUnprocessed().RunAndReturn(func() ([]models.Review, error) {
		close(entered)
		<-release
		return nil, nil
	})

	done := make(chan error, 1)
	go func() { done <- proc.Run(context.Background()) }()

	<-entered // first Run is inside GetUnprocessed; isRunning == true

	err := proc.Run(context.Background()) // second call: skipped
	require.NoError(t, err)

	close(release)
	require.NoError(t, <-done)
}
