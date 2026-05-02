package handlers

import (
	"errors"
	"majestic-gondola/internal/apperr"
	"majestic-gondola/internal/dto"
	"net/http"

	"github.com/gin-gonic/gin"
)

func respondErr(c *gin.Context, err error) {
	if appErr, ok := errors.AsType[*apperr.AppError](err); ok {
		c.JSON(appErr.Code, dto.ErrResponse{Message: appErr.Message})
		return
	}
	c.JSON(http.StatusInternalServerError, dto.ErrResponse{Message: "Internal server error"})
}
