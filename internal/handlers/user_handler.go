package handlers

import (
	"log/slog"
	"majestic-gondola/internal/dto"
	"majestic-gondola/internal/mappers"
	"majestic-gondola/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	log         *slog.Logger
	userService service.UserService
}

func NewUserHandler(userService service.UserService, logger *slog.Logger) *UserHandler {
	return &UserHandler{log: logger.With("component", "user_handler"), userService: userService}
}

// GetUsers godoc
//
//	@Summary		List users
//	@Description	Get a list of all users in the database
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Success		200	{array}		dto.UserResponse
//	@Failure		500	{object}	dto.ErrResponse	"Internal server error"
//	@Router			/users [get]
func (h *UserHandler) GetUsers(c *gin.Context) {
	users, err := h.userService.GetAll()
	if err != nil {
		respondErr(c, err)
		return
	}

	c.JSON(http.StatusOK, mappers.ToUserResponseList(users))
}

// GetUser godoc
//
//	@Summary		Get a user
//	@Description	Retrieve a single user by their unique ID
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"User ID"
//	@Success		200	{object}	dto.UserResponse
//	@Failure		400	{object}	dto.ErrResponse	"Invalid ID format"
//	@Failure		404	{object}	dto.ErrResponse	"User not found"
//	@Failure		500	{object}	dto.ErrResponse	"Internal server error"
//	@Router			/users/{id} [get]
func (h *UserHandler) GetUser(c *gin.Context) {
	var req dto.IdUriRequest

	if err := c.ShouldBindUri(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrResponse{Message: err.Error()})
		return
	}

	user, err := h.userService.Get(req.Id)

	if err != nil {
		respondErr(c, err)
		return
	}

	c.JSON(http.StatusOK, mappers.ToUserResponse(user))
}

// CreateUsers godoc
//
//	@Summary		Bulk create users
//	@Description	Create multiple users at once from a JSON array
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			users	body	[]dto.CreateUserRequest	true	"List of users to create"
//	@Success		201		"Created"
//	@Failure		400		{object}	dto.ErrResponse	"Invalid request body"
//	@Failure		500		{object}	dto.ErrResponse	"Internal server error"
//	@Router			/users [post]
func (h *UserHandler) CreateUsers(c *gin.Context) {
	var req []dto.CreateUserRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrResponse{Message: err.Error()})
		return
	}

	users := mappers.CreateToUserList(req)
	err := h.userService.BulkCreate(users)

	if err != nil {
		respondErr(c, err)
		return
	}

	c.Status(http.StatusCreated)
}

// UpdateUser godoc
//
//	@Summary		Update a user
//	@Description	Update the details of an existing user
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			id		path	int						true	"User ID"
//	@Param			user	body	dto.UpdateUserRequest	true	"User update data"
//	@Success		200		"Updated"
//	@Failure		400		{object}	dto.ErrResponse	"Invalid request body"
//	@Failure		500		{object}	dto.ErrResponse	"Internal server error"
//	@Router			/users/{id} [put]
func (h *UserHandler) UpdateUser(c *gin.Context) {
	var uri dto.IdUriRequest
	var body dto.UpdateUserRequest

	if err := c.ShouldBindUri(&uri); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrResponse{Message: err.Error()})
		return
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrResponse{Message: err.Error()})
		return
	}

	user := mappers.UpdateToUser(uri.Id, body)

	err := h.userService.Update(user)
	if err != nil {
		respondErr(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}
