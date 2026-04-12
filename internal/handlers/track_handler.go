package handlers

import (
	"fmt"
	"log/slog"
	"majestic-gondola/internal/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-pg/pg/v10"
)

type TrackHandler struct {
	db     *pg.DB
	logger *slog.Logger
}

func NewTrackHandler(db *pg.DB, logger *slog.Logger) *TrackHandler {
	return &TrackHandler{db: db, logger: logger}
}

// GetTracks godoc
//
//	@Summary		Get tracks
//	@Description	Get all tracks or filter by ID via query parameter
//	@Tags			tracks
//	@Accept			json
//	@Produce		json
//	@Param			id	query	int	false	"Track ID"
//	@Success		200
//	@Failure		500
//	@Router			/track [get]
func (h *TrackHandler) GetTracks(c *gin.Context) {
	var tracks []models.Track
	id := c.Query("id")

	if len(id) > 0 {
		h.db.Model(&tracks).Where("id = ?", id).Select()
	} else {
		h.db.Model(&tracks).Select()
	}

	c.JSON(http.StatusOK, tracks)
}

// CreateTracks godoc
//
//	@Summary		Create new tracks
//	@Description	Bulk create tracks from a JSON array
//	@Tags			tracks
//	@Accept			json
//	@Produce		json
//	@Param			tracks	body	[]models.CreateTrackRequest	true	"Array of tracks to create"
//	@Success		200
//	@Failure		400
//	@Router			/track [post]
func (h *TrackHandler) CreateTracks(c *gin.Context) {
	var req []models.CreateTrackRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tracks := make([]*models.Track, 0, len(req))

	for _, reqTrack := range req {
		track := models.Track{
			Name:        reqTrack.Name,
			Author:      reqTrack.Author,
			ReleaseDate: reqTrack.ReleaseDate,
			Genres:      reqTrack.Genres,
		}

		tracks = append(tracks, &track)
	}

	// TODO: Wrap model creation with logger?
	_, err := h.db.Model(&tracks).Insert()

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	h.logger.Info("Created a new track")
}

// PopulateTracks godoc
//
//	@Summary		Populate dummy tracks
//	@Description	Generate a specified number of dummy tracks for testing
//	@Tags			tracks
//	@Accept			json
//	@Produce		json
//	@Param			count	path	int	true	"Number of tracks to generate"
//	@Success		200
//	@Failure		400
//	@Router			/track/populate/{count} [post]
func (h *TrackHandler) PopulateTracks(c *gin.Context) {

	req := struct {
		Count int `uri:"count" binding:"required"`
	}{}

	if err := c.ShouldBindUri(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tracks := make([]*models.Track, 0, req.Count)

	for i := 0; i < req.Count; i++ {
		track := models.Track{
			Name:        fmt.Sprintf("Track %d", i),
			Author:      fmt.Sprintf("Author %d", i),
			ReleaseDate: time.Now(),
			Genres:      []string{"Tag"},
		}

		tracks = append(tracks, &track)
	}

	// TODO: Wrap model creation with logger?
	_, err := h.db.Model(&tracks).Insert()

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	h.logger.Info("Created several new tracks")
}
