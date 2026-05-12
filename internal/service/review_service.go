package service

import (
	"errors"
	"log/slog"
	"majestic-gondola/internal/apperr"
	"majestic-gondola/internal/models"
	"majestic-gondola/internal/repository"
)

type ReviewService interface {
	Get(id int) (*models.Review, error)
	GetAll() ([]models.Review, error)
	BulkCreate(reviews []*models.Review) error
	Update(review *models.Review) error
	GetUnprocessed() ([]models.Review, error)
	BulkDelete(ids []int) error
	GetTrackReviews(trackId int) ([]models.Review, error)
	GetUserReviews(userId int) ([]models.Review, error)
}

type reviewService struct {
	reviewRepository repository.ReviewRepository
	log              *slog.Logger
}

func NewReviewService(reviewRepository repository.ReviewRepository, logger *slog.Logger) ReviewService {
	return &reviewService{log: logger.With("component", "review_service"), reviewRepository: reviewRepository}
}

// BulkCreate implements [ReviewService].
func (r *reviewService) BulkCreate(reviews []*models.Review) error {
	err := r.reviewRepository.BulkCreate(reviews)

	if err != nil {
		r.log.Error("Failed to bulk create reviews", slog.Any("error", err), slog.Int("parsed_reviews", len(reviews)))
		return apperr.Internal(err)
	}

	r.log.Info("Bulk created reviews", slog.Int("count", len(reviews)))
	return nil
}

// Get implements [ReviewService].
func (r *reviewService) Get(id int) (*models.Review, error) {
	review, err := r.reviewRepository.FindById(id)

	if err != nil {
		if errors.Is(err, apperr.ErrNotFound) {
			return nil, apperr.NotFound("Review not found", err)
		}
		r.log.Error("Failed to fetch a review by id", slog.Any("error", err), slog.Int("id", id))
		return nil, apperr.Internal(err)
	}

	return review, nil
}

// GetAll implements [ReviewService].
func (r *reviewService) GetAll() ([]models.Review, error) {
	reviews, err := r.reviewRepository.GetAll()

	if err != nil {
		r.log.Error("Failed to fetch reviews", slog.Any("error", err))
		return nil, apperr.Internal(err)
	}

	return reviews, nil
}

// Update implements [ReviewService].
func (r *reviewService) Update(review *models.Review) error {
	err := r.reviewRepository.Update(review)

	if err != nil {
		if errors.Is(err, apperr.ErrNotFound) {
			return apperr.NotFound("Review not found", err)
		}
		r.log.Error("Failed to update the review", slog.Any("error", err), slog.Any("parsed_review", review))
		return apperr.Internal(err)
	}

	r.log.Info("Updated a review")
	return nil
}

// BulkDelete implements [ReviewService].
func (r *reviewService) BulkDelete(ids []int) error {
	err := r.reviewRepository.BulkDelete(ids)

	if err != nil {
		if errors.Is(err, apperr.ErrNotFound) {
			return apperr.NotFound("Review not found", err)
		}
		r.log.Error("Failed to bulk delete reviews", slog.Any("error", err), slog.Int("ids_length", len(ids)))
		return apperr.Internal(err)
	}

	r.log.Info("Bulk deleted reviews", slog.Int("count", len(ids)))
	return nil
}

// GetUnprocessed implements [ReviewService].
func (r *reviewService) GetUnprocessed() ([]models.Review, error) {
	review, err := r.reviewRepository.GetUnprocessed()

	if err != nil {
		if errors.Is(err, apperr.ErrNotFound) {
			return nil, apperr.NotFound("Unprocessed reviews not found", err)
		}
		r.log.Error("Failed to fetch unprocessed reviews", slog.Any("error", err))
		return nil, apperr.Internal(err)
	}

	return review, nil
}

// GetTrackReviews implements [ReviewService].
func (r *reviewService) GetTrackReviews(trackId int) ([]models.Review, error) {
	review, err := r.reviewRepository.GetTrackReviews(trackId)

	if err != nil {
		if errors.Is(err, apperr.ErrNotFound) {
			return nil, apperr.NotFound("Reviews not found", err)
		}
		r.log.Error("Failed to fetch reviews for track", slog.Any("error", err), slog.Int("track_id", trackId))
		return nil, apperr.Internal(err)
	}

	return review, nil
}

// GetUserReviews implements [ReviewService].
func (r *reviewService) GetUserReviews(userId int) ([]models.Review, error) {
	review, err := r.reviewRepository.GetUserReviews(userId)

	if err != nil {
		if errors.Is(err, apperr.ErrNotFound) {
			return nil, apperr.NotFound("Reviews not found", err)
		}
		r.log.Error("Failed to fetch reviews for user", slog.Any("error", err), slog.Int("user_id", userId))
		return nil, apperr.Internal(err)
	}

	return review, nil
}
