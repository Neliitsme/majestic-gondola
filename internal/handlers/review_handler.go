package handlers

import (
	"log/slog"
	"majestic-gondola/internal/dto"
	"majestic-gondola/internal/mappers"
	"majestic-gondola/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ReviewHandler struct {
	log           *slog.Logger
	reviewService service.ReviewService
}

func NewReviewHandler(reviewService service.ReviewService, logger *slog.Logger) *ReviewHandler {
	return &ReviewHandler{log: logger.With("component", "review_handler"), reviewService: reviewService}
}

// GetReviews godoc
//
//	@Summary		List reviews
//	@Description	Get a list of all reviews in the database
//	@Tags			reviews
//	@Accept			json
//	@Produce		json
//	@Success		200	{array}		dto.ReviewResponse
//	@Failure		500	{object}	dto.ErrResponse	"Internal server error"
//	@Router			/reviews [get]
func (h *ReviewHandler) GetReviews(c *gin.Context) {
	reviews, err := h.reviewService.GetAll()
	if err != nil {
		respondErr(c, err)
		return
	}

	c.JSON(http.StatusOK, mappers.ToReviewResponseList(reviews))
}

// GetReview godoc
//
//	@Summary		Get a review
//	@Description	Retrieve a single review by its unique ID
//	@Tags			reviews
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"Review ID"
//	@Success		200	{object}	dto.ReviewResponse
//	@Failure		400	{object}	dto.ErrResponse	"Invalid ID format"
//	@Failure		404	{object}	dto.ErrResponse	"Review not found"
//	@Failure		500	{object}	dto.ErrResponse	"Internal server error"
//	@Router			/reviews/{id} [get]
func (h *ReviewHandler) GetReview(c *gin.Context) {
	var req dto.IdUriRequest

	if err := c.ShouldBindUri(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrResponse{Message: err.Error()})
		return
	}

	review, err := h.reviewService.Get(req.Id)

	if err != nil {
		respondErr(c, err)
		return
	}

	c.JSON(http.StatusOK, mappers.ToReviewResponse(review))
}

// CreateReviews godoc
//
//	@Summary		Bulk create reviews
//	@Description	Create multiple reviews at once from a JSON array
//	@Tags			reviews
//	@Accept			json
//	@Produce		json
//	@Param			reviews	body	[]dto.CreateReviewRequest	true	"List of reviews to create"
//	@Success		201		"Created"
//	@Failure		400		{object}	dto.ErrResponse	"Invalid request body"
//	@Failure		500		{object}	dto.ErrResponse	"Internal server error"
//	@Router			/reviews [post]
func (h *ReviewHandler) CreateReviews(c *gin.Context) {
	var req []dto.CreateReviewRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrResponse{Message: err.Error()})
		return
	}

	reviews := mappers.CreateToReviewList(req)
	err := h.reviewService.BulkCreate(reviews)

	if err != nil {
		respondErr(c, err)
		return
	}

	c.Status(http.StatusCreated)
}

// UpdateReview godoc
//
//	@Summary		Update a review
//	@Description	Update the details of an existing review
//	@Tags			reviews
//	@Accept			json
//	@Produce		json
//	@Param			id		path	int						true	"Review ID"
//	@Param			review	body	dto.UpdateReviewRequest	true	"Review update data"
//	@Success		200		"Updated"
//	@Failure		400		{object}	dto.ErrResponse	"Invalid request body"
//	@Failure		500		{object}	dto.ErrResponse	"Internal server error"
//	@Router			/reviews/{id} [put]
func (h *ReviewHandler) UpdateReview(c *gin.Context) {
	var uri dto.IdUriRequest
	var body dto.UpdateReviewRequest

	if err := c.ShouldBindUri(&uri); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrResponse{Message: err.Error()})
		return
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrResponse{Message: err.Error()})
		return
	}

	review := mappers.UpdateToReview(uri.Id, body)

	err := h.reviewService.Update(review)
	if err != nil {
		respondErr(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}

// GetTrackReviews godoc
//
//	@Summary		List track reviews
//	@Description	Get a list of all track reviews in the database
//	@Tags			reviews
//	@Accept			json
//	@Produce		json
//	@Param			id	path	int	true	"Track ID"
//	@Success		200
//	@Failure		500	{object}	dto.ErrResponse	"Internal server error"
//	@Router			/tracks/{id}/reviews [get]
func (h *ReviewHandler) GetTrackReviews(c *gin.Context) {
	var req dto.IdUriRequest

	if err := c.ShouldBindUri(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrResponse{Message: err.Error()})
		return
	}

	reviews, err := h.reviewService.GetTrackReviews(req.Id)

	if err != nil {
		respondErr(c, err)
		return
	}

	c.JSON(http.StatusOK, mappers.ToReviewResponseList(reviews))
}

// GetUserReviews godoc
//
//	@Summary		List user reviews
//	@Description	Get a list of all reviews of the user in the database
//	@Tags			reviews
//	@Accept			json
//	@Produce		json
//	@Param			id	path	int	true	"User ID"
//	@Success		200
//	@Failure		500	{object}	dto.ErrResponse	"Internal server error"
//	@Router			/users/{id}/reviews [get]
func (h *ReviewHandler) GetUserReviews(c *gin.Context) {
	var req dto.IdUriRequest

	if err := c.ShouldBindUri(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrResponse{Message: err.Error()})
		return
	}

	reviews, err := h.reviewService.GetUserReviews(req.Id)

	if err != nil {
		respondErr(c, err)
		return
	}

	c.JSON(http.StatusOK, mappers.ToReviewResponseList(reviews))
}
