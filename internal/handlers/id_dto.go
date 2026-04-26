package handlers

type IdUriRequest struct {
	Id int `uri:"id" binding:"required"`
}
