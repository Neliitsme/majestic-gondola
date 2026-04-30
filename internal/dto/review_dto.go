package dto

type CreateReviewRequest struct {
	UserId  *int `json:"user_id" example:"1"`
	TrackId *int `json:"track_id" example:"1"`
	Score   int  `json:"score" binding:"required,min=0,max=100" example:"1"`
}

type UpdateReviewRequest struct {
	UserId  *int `json:"user_id" example:"1"`
	TrackId *int `json:"track_id" example:"1"`
	Score   int  `json:"score" binding:"required,min=0,max=100" example:"1"`
}

type ReviewResponse struct {
	Id          int    `json:"id" example:"1"`
	UserId      *int   `json:"user_id,omitempty" example:"1"`
	TrackId     *int   `json:"track_id,omitempty" example:"1"`
	Score       int    `json:"score" example:"1"`
	IsDeleted   bool   `json:"is_deleted" example:"false"`
	IsProcessed bool   `json:"is_processed" example:"false"`
	CreatedAt   string `json:"created_at" example:"2006-01-02 15:04:05"`
}
