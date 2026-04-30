package dto

type CreateUserRequest struct {
	Name string `json:"name" binding:"required" example:"Tempo"`
}

type UpdateUserRequest struct {
	Name string `json:"name" binding:"required" example:"Tempo"`
}

type UserResponse struct {
	Id        int    `json:"id" example:"1"`
	Name      string `json:"name" example:"Tempo"`
	CreatedAt string `json:"created_at" example:"2006-01-02 15:04:05"`
}
