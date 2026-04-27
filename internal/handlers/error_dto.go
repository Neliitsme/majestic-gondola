package handlers

type ErrResponse struct {
	Message string `json:"message" example:"error"`
}

var ErrInternalMsg = "Internal server error"
