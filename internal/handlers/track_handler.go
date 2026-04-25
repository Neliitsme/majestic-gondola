package handlers

import (
	"log/slog"
	"majestic-gondola/internal/models"
	"majestic-gondola/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-pg/pg/v10"
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
//	@Success		200	{array}		TrackResponse
//	@Failure		500	{object}	map[string]string	"Internal server error"
//	@Router			/track [get]
func (h *TrackHandler) GetTracks(c *gin.Context) {
	tracks, err := h.trackService.GetAll()
	if err != nil {
		h.log.Error("Failed to fetch the track list", slog.Any("error", err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	resTracks := make([]TrackResponse, 0, len(tracks))
	for i := range tracks {
		resTracks = append(resTracks, ToTrackResponse(&tracks[i]))
	}

	c.JSON(http.StatusOK, resTracks)
}

// GetTrack godoc
//
//	@Summary		Get a track
//	@Description	Retrieve a single track by its unique ID
//	@Tags			tracks
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"Track ID"
//	@Success		200	{object}	TrackResponse
//	@Failure		400	{object}	map[string]string	"Invalid ID format"
//	@Failure		404	{object}	map[string]string	"Track not found"
//	@Failure		500	{object}	map[string]string	"Internal server error"
//	@Router			/track/{id} [get]
func (h *TrackHandler) GetTrack(c *gin.Context) {
	var req IdUriRequest

	if err := c.ShouldBindUri(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	track, err := h.trackService.Get(req.Id)

	if err == pg.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "Track not found"})
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	c.JSON(http.StatusOK, ToTrackResponse(track))
}

// CreateTracks godoc
//
//	@Summary		Bulk create tracks
//	@Description	Create multiple tracks at once from a JSON array
//	@Tags			tracks
//	@Accept			json
//	@Produce		json
//	@Param			tracks	body	[]CreateTrackRequest	true	"List of tracks to create"
//	@Success		201		"Created"
//	@Failure		400		{object}	map[string]string	"Invalid request body"
//	@Failure		500		{object}	map[string]string	"Internal server error"
//	@Router			/track [post]
func (h *TrackHandler) CreateTracks(c *gin.Context) {
	var req []CreateTrackRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tracks := make([]*models.Track, 0, len(req))
	// TODO: Make
	for i := range req {
		track, err := CreateToTrack(req[i])
		if err != nil {
			h.log.Error("Failed to map Create to Track", slog.Any("error", err), slog.Any("track_request", req[i]))
			c.JSON(http.StatusBadRequest, gin.H{"error": "Error while trying to parse the request"})
			return
		}
		tracks = append(tracks, track)
	}

	err := h.trackService.BulkCreate(tracks)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
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
//	@Param			count	path	int	true	"Number of tracks to generate"
//	@Success		201		"Created"
//	@Failure		400		{object}	map[string]string	"Invalid count"
//	@Failure		500		{object}	map[string]string	"Internal server error"
//	@Router			/track/populate/{count} [post]
func (h *TrackHandler) PopulateTracks(c *gin.Context) {
	req := struct {
		Count int `uri:"count" binding:"required"`
	}{}

	if err := c.ShouldBindUri(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.trackService.Generate(req.Count)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
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
//	@Param			id		path	int					true	"Track ID"
//	@Param			track	body	UpdateTrackRequest	true	"Track update data"
//	@Success		200		"Updated"
//	@Failure		400		{object}	map[string]string	"Invalid request body"
//	@Failure		500		{object}	map[string]string	"Internal server error"
//	@Router			/track/{id} [put]
func (h *TrackHandler) UpdateTrack(c *gin.Context) {
	var uri IdUriRequest
	var body UpdateTrackRequest

	if err := c.ShouldBindUri(&uri); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	track, err := UpdateToTrack(uri.Id, body)
	if err != nil {
		h.log.Error("Failed to map Update to Track", slog.Any("error", err), slog.Int("track_id", uri.Id), slog.Any("track_request", body))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error while trying to parse the request"})
		return
	}

	err = h.trackService.Update(track)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	c.Status(http.StatusOK)
}
