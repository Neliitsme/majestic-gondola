package mappers

import (
	"majestic-gondola/internal/dto"
	"majestic-gondola/internal/models"
	"time"
)

func ToReviewResponse(r *models.Review) dto.ReviewResponse {
	return dto.ReviewResponse{
		Id:          r.Id,
		UserId:      r.UserId,
		TrackId:     r.TrackId,
		Score:       r.Score,
		IsDeleted:   r.IsDeleted,
		IsProcessed: r.IsProcessed,
		CreatedAt:   r.CreatedAt.Format(time.DateTime),
	}
}

func ToReviewResponseList(rs []models.Review) []dto.ReviewResponse {
	responses := make([]dto.ReviewResponse, 0, len(rs))
	for i := range rs {
		responses = append(responses, ToReviewResponse(&rs[i]))
	}
	return responses
}

func CreateToReview(rr dto.CreateReviewRequest) *models.Review {
	return &models.Review{
		UserId:  rr.UserId,
		TrackId: rr.TrackId,
		Score:   rr.Score,
	}
}

func CreateToReviewList(rrs []dto.CreateReviewRequest) []*models.Review {
	reviews := make([]*models.Review, 0, len(rrs))
	for i := range rrs {
		review := CreateToReview(rrs[i])
		reviews = append(reviews, review)
	}
	return reviews
}

func UpdateToReview(id int, rr dto.UpdateReviewRequest) *models.Review {
	return &models.Review{
		Id:      id,
		UserId:  rr.UserId,
		TrackId: rr.TrackId,
		Score:   rr.Score,
	}
}
