package handlers

import (
	"errors"
	"log/slog"
	"majestic-gondola/internal/apperr"
	"majestic-gondola/internal/dto"
	"majestic-gondola/internal/mappers"
	"majestic-gondola/internal/models"
	"majestic-gondola/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type TrackHandler struct {
	log          *slog.Logger
	trackService service.TrackService
}

func NewTrackHandler(trackService service.TrackService, logger *slog.Logger) *TrackHandler {
	return &TrackHandler{log: logger.With("component", "track_handler"), trackService: trackService}
}

// GetTracks godoc
//
//	@Summary		List tracks
//	@Description	Get a list of all tracks in the database
//	@Tags			tracks
//	@Accept			json
//	@Produce		json
//	@Success		200	{array}		dto.TrackResponse
//	@Failure		500	{object}	dto.ErrResponse	"Internal server error"
//	@Router			/tracks [get]
func (h *TrackHandler) GetTracks(c *gin.Context) {
	tracks, err := h.trackService.GetAll()
	if err != nil {
		h.log.Error("Failed to fetch the track list", slog.Any("error", err))
		c.JSON(http.StatusInternalServerError, dto.InternalErrResponse)
		return
	}

	c.JSON(http.StatusOK, mappers.ToTrackResponseList(tracks))
}

// GetTrack godoc
//
//	@Summary		Get a track
//	@Description	Retrieve a single track by its unique ID
//	@Tags			tracks
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"Track ID"
//	@Success		200	{object}	dto.TrackResponse
//	@Failure		400	{object}	dto.ErrResponse	"Invalid ID format"
//	@Failure		404	{object}	dto.ErrResponse	"Track not found"
//	@Failure		500	{object}	dto.ErrResponse	"Internal server error"
//	@Router			/tracks/{id} [get]
func (h *TrackHandler) GetTrack(c *gin.Context) {
	var req dto.IdUriRequest

	if err := c.ShouldBindUri(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrResponse{Message: err.Error()})
		return
	}

	track, err := h.trackService.Get(req.Id)

	if err != nil {
		if appErr, ok := errors.AsType[*apperr.AppError](err); ok {
			c.JSON(appErr.Code, dto.ErrResponse{Message: appErr.Message})
			return
		}
		c.JSON(http.StatusInternalServerError, dto.InternalErrResponse)
		return
	}

	c.JSON(http.StatusOK, mappers.ToTrackResponse(track))
}

// CreateTracks godoc
//
//	@Summary		Bulk create tracks
//	@Description	Create multiple tracks at once from a JSON array
//	@Tags			tracks
//	@Accept			json
//	@Produce		json
//	@Param			tracks	body	[]dto.CreateTrackRequest	true	"List of tracks to create"
//	@Success		201		"Created"
//	@Failure		400		{object}	dto.ErrResponse	"Invalid request body"
//	@Failure		500		{object}	dto.ErrResponse	"Internal server error"
//	@Router			/tracks [post]
func (h *TrackHandler) CreateTracks(c *gin.Context) {
	var req []dto.CreateTrackRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrResponse{Message: err.Error()})
		return
	}

	tracks := make([]*models.Track, 0, len(req))
	for i := range req {
		track, err := mappers.CreateToTrack(req[i])
		if err != nil {
			h.log.Error("Failed to map Create to Track", slog.Any("error", err), slog.Any("track_request", req[i]))
			c.JSON(http.StatusBadRequest, dto.ErrResponse{Message: "Error while trying to parse the request"})
			return
		}
		tracks = append(tracks, track)
	}

	err := h.trackService.BulkCreate(tracks)

	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.InternalErrResponse)
		return
	}

	c.Status(http.StatusCreated)
}

// PopulateTracks godoc
//
//	@Summary		Seed dummy data
//	@Description	Generate random tracks for development testing
//	@Tags			dev
//	@Accept			json
//	@Produce		json
//	@Param			count	path	int					true	"Number of tracks to generate"
//	@Param			track	body	dto.PopulateRequest	false	"Populate data"
//	@Success		201		"Created"
//	@Failure		400		{object}	dto.ErrResponse	"Invalid count"
//	@Failure		500		{object}	dto.ErrResponse	"Internal server error"
//	@Router			/tracks/populate/{count} [post]
func (h *TrackHandler) PopulateTracks(c *gin.Context) {
	uri := struct {
		Count int `uri:"count" binding:"required"`
	}{}
	var body dto.PopulateRequest

	if err := c.ShouldBindUri(&uri); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrResponse{Message: err.Error()})
		return
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrResponse{Message: err.Error()})
		return
	}

	err := h.trackService.Generate(uri.Count, body.ArtistId)

	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.InternalErrResponse)
		return
	}

	c.Status(http.StatusCreated)
}

// UpdateTrack godoc
//
//	@Summary		Update a track
//	@Description	Update the details of an existing track
//	@Tags			tracks
//	@Accept			json
//	@Produce		json
//	@Param			id		path	int						true	"Track ID"
//	@Param			track	body	dto.UpdateTrackRequest	true	"Track update data"
//	@Success		200		"Updated"
//	@Failure		400		{object}	dto.ErrResponse	"Invalid request body"
//	@Failure		500		{object}	dto.ErrResponse	"Internal server error"
//	@Router			/tracks/{id} [put]
func (h *TrackHandler) UpdateTrack(c *gin.Context) {
	var uri dto.IdUriRequest
	var body dto.UpdateTrackRequest

	if err := c.ShouldBindUri(&uri); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrResponse{Message: err.Error()})
		return
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrResponse{Message: err.Error()})
		return
	}

	track, err := mappers.UpdateToTrack(uri.Id, body)
	if err != nil {
		h.log.Error("Failed to map Update to Track", slog.Any("error", err), slog.Int("track_id", uri.Id), slog.Any("track_request", body))
		c.JSON(http.StatusBadRequest, dto.ErrResponse{Message: "Error while trying to parse the request"})
		return
	}

	err = h.trackService.Update(track)
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
