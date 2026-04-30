package handlers

import (
	"errors"
	"log/slog"
	"majestic-gondola/internal/apperr"
	"majestic-gondola/internal/dto"
	"majestic-gondola/internal/mappers"
	"majestic-gondola/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ArtistHandler struct {
	log           *slog.Logger
	artistService service.ArtistService
}

func NewArtistHandler(artistService service.ArtistService, logger *slog.Logger) *ArtistHandler {
	return &ArtistHandler{log: logger.With("component", "artist_handler"), artistService: artistService}
}

// GetArtists godoc
//
//	@Summary		List artists
//	@Description	Get a list of all artists in the database
//	@Tags			artists
//	@Accept			json
//	@Produce		json
//	@Success		200	{array}		dto.ArtistResponse
//	@Failure		500	{object}	dto.ErrResponse	"Internal server error"
//	@Router			/artists [get]
func (h *ArtistHandler) GetArtists(c *gin.Context) {
	artists, err := h.artistService.GetAll()
	if err != nil {
		h.log.Error("Failed to fetch the artist list", slog.Any("error", err))
		c.JSON(http.StatusInternalServerError, dto.InternalErrResponse)
		return
	}

	c.JSON(http.StatusOK, mappers.ToArtistResponseList(artists))
}

// GetArtist godoc
//
//	@Summary		Get a artist
//	@Description	Retrieve a single artist by their unique ID
//	@Tags			artists
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"Artist ID"
//	@Success		200	{object}	dto.ArtistResponse
//	@Failure		400	{object}	dto.ErrResponse	"Invalid ID format"
//	@Failure		404	{object}	dto.ErrResponse	"Artist not found"
//	@Failure		500	{object}	dto.ErrResponse	"Internal server error"
//	@Router			/artists/{id} [get]
func (h *ArtistHandler) GetArtist(c *gin.Context) {
	var req dto.IdUriRequest

	if err := c.ShouldBindUri(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrResponse{Message: err.Error()})
		return
	}

	artist, err := h.artistService.Get(req.Id)

	if err != nil {
		if appErr, ok := errors.AsType[*apperr.AppError](err); ok {
			c.JSON(appErr.Code, dto.ErrResponse{Message: appErr.Message})
			return
		}
		c.JSON(http.StatusInternalServerError, dto.InternalErrResponse)
		return
	}

	c.JSON(http.StatusOK, mappers.ToArtistResponse(artist))
}

// CreateArtists godoc
//
//	@Summary		Bulk create artists
//	@Description	Create multiple artists at once from a JSON array
//	@Tags			artists
//	@Accept			json
//	@Produce		json
//	@Param			artists	body	[]dto.CreateArtistRequest	true	"List of artists to create"
//	@Success		201		"Created"
//	@Failure		400		{object}	dto.ErrResponse	"Invalid request body"
//	@Failure		500		{object}	dto.ErrResponse	"Internal server error"
//	@Router			/artists [post]
func (h *ArtistHandler) CreateArtists(c *gin.Context) {
	var req []dto.CreateArtistRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrResponse{Message: err.Error()})
		return
	}

	artists := mappers.CreateToArtistList(req)
	err := h.artistService.BulkCreate(artists)

	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.InternalErrResponse)
		return
	}

	c.Status(http.StatusCreated)
}

// UpdateArtist godoc
//
//	@Summary		Update an artist
//	@Description	Update the details of an existing artist
//	@Tags			artists
//	@Accept			json
//	@Produce		json
//	@Param			id		path	int						true	"Artist ID"
//	@Param			artist	body	dto.UpdateArtistRequest	true	"Artist update data"
//	@Success		200		"Updated"
//	@Failure		400		{object}	dto.ErrResponse	"Invalid request body"
//	@Failure		500		{object}	dto.ErrResponse	"Internal server error"
//	@Router			/artists/{id} [put]
func (h *ArtistHandler) UpdateArtist(c *gin.Context) {
	var uri dto.IdUriRequest
	var body dto.UpdateArtistRequest

	if err := c.ShouldBindUri(&uri); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrResponse{Message: err.Error()})
		return
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrResponse{Message: err.Error()})
		return
	}

	artist := mappers.UpdateToArtist(uri.Id, body)

	err := h.artistService.Update(artist)
	if err != nil {
		if appErr, ok := errors.AsType[*apperr.AppError](err); ok {
			c.JSON(appErr.Code, dto.ErrResponse{Message: appErr.Message})
			return
		}
		c.JSON(http.StatusInternalServerError, dto.InternalErrResponse)
		return
	}

	c.Status(http.StatusNoContent)
}
