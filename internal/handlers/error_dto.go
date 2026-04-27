package handlers

type ErrResponse struct {
	Message string `json:"message" example:"error"`
}

var InternalErrResponse = ErrResponse{Message: "Internal server error"}
