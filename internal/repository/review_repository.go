package repository

import (
	"errors"
	"log/slog"
	"majestic-gondola/internal/apperr"
	"majestic-gondola/internal/models"

	"github.com/go-pg/pg/v10"
)

type ReviewRepository interface {
	FindById(id int) (*models.Review, error)
	GetAll() ([]models.Review, error)
	BulkCreate(reviews []*models.Review) error
	Update(review *models.Review) error
	GetUnprocessed() ([]models.Review, error)
	BulkDelete(ids []int) error
	GetTrackReviews(trackId int) ([]models.Review, error)
	GetUserReviews(userId int) ([]models.Review, error)
}

type reviewRepository struct {
	db  *pg.DB
	log *slog.Logger
}

func NewReviewRepository(db *pg.DB, logger *slog.Logger) ReviewRepository {
	return &reviewRepository{db: db, log: logger.With("component", "review_repository")}
}

// BulkCreate implements [ReviewRepository].
func (r *reviewRepository) BulkCreate(reviews []*models.Review) error {
	_, err := r.db.Model(&reviews).Insert()

	if err != nil {
		return err
	}

	r.log.Info("Finished creating reviews")
	return nil
}

// FindById implements [ReviewRepository].
func (r *reviewRepository) FindById(id int) (*models.Review, error) {
	review := new(models.Review)

	err := r.db.Model(review).Where("review_id = ? AND is_deleted = false", id).Select()

	if err != nil {
		if errors.Is(err, pg.ErrNoRows) {
			return nil, apperr.ErrNotFound
		}
		return nil, err
	}

	return review, nil
}

// GetAll implements [ReviewRepository].
func (r *reviewRepository) GetAll() ([]models.Review, error) {
	var reviews []models.Review
	err := r.db.Model(&reviews).Where("is_deleted = false").Select()

	if err != nil {
		return nil, err
	}

	return reviews, nil
}

// Update implements [ReviewRepository].
func (r *reviewRepository) Update(review *models.Review) error {
	res, err := r.db.Model(review).ExcludeColumn("created_at").WherePK().Update()

	if err != nil {
		return err
	}

	if res.RowsAffected() == 0 {
		return apperr.ErrNotFound
	}

	r.log.Info("Finished updating review")
	return nil
}

// GetUnprocessed implements [ReviewRepository].
func (r *reviewRepository) GetUnprocessed() ([]models.Review, error) {
	var reviews []models.Review
	err := r.db.Model(&reviews).
		Relation("Track").
		Where("review.is_processed = false AND review.is_deleted = false").
		Select()

	if err != nil {
		return nil, err
	}

	return reviews, nil
}

// BulkDelete implements [ReviewRepository].
func (r *reviewRepository) BulkDelete(ids []int) error {
	res, err := r.db.Model(new(models.Review)).
		Set("is_deleted = true").
		Where("review_id IN (?) AND is_deleted = false", pg.In(ids)).
		Update()

	if err != nil {
		return err
	}

	if res.RowsAffected() == 0 {
		return apperr.ErrNotFound
	}

	r.log.Info("Finished deleting reviews")
	return nil
}

// GetTrackReviews implements [ReviewRepository].
func (r *reviewRepository) GetTrackReviews(trackId int) ([]models.Review, error) {
	var reviews []models.Review
	err := r.db.Model(&reviews).Where("track_id = ? AND is_deleted = false", trackId).Select()

	if err != nil {
		return nil, err
	}

	return reviews, nil
}

// GetUserReviews implements [ReviewRepository].
func (r *reviewRepository) GetUserReviews(userId int) ([]models.Review, error) {
	var reviews []models.Review
	err := r.db.Model(&reviews).Where("user_id = ? AND is_deleted = false", userId).Select()

	if err != nil {
		return nil, err
	}

	return reviews, nil
}
