package handlers

import (
	"fmt"
	"log/slog"
	"majestic-gondola/internal/models"
	"majestic-gondola/internal/repository"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type TrackHandler struct {
	log             *slog.Logger
	trackRepository repository.TrackStore
}

func NewTrackHandler(trackRepository *repository.TrackRepository, logger *slog.Logger) *TrackHandler {
	return &TrackHandler{log: logger.With("component", "track_handler"), trackRepository: trackRepository}
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
	strId := c.Query("id")

	if len(strId) == 0 {
		rTracks, err := h.trackRepository.GetAll()
		if err != nil {
			panic(err)
		}

		tracks = rTracks
	} else if id, err := strconv.Atoi(strId); err == nil {
		track, err := h.trackRepository.FindById(id)
		if err != nil {
			panic(err)
		}

		tracks = []models.Track{*track}
	} else {
		c.JSON(http.StatusBadRequest, "Bad query param")
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

	err := h.trackRepository.BulkCreate(tracks)

	if err != nil {
		panic(err)
	}

	h.log.Info("Created a new track")
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

	err := h.trackRepository.BulkCreate(tracks)

	if err != nil {
		panic(err)
	}

	h.log.Info("Created several new tracks")
}

// UpdateTrack godoc
//
//	@Summary		Update an existing track
//	@Description	Update track details based on the provided JSON body
//	@Tags			tracks
//	@Accept			json
//	@Produce		json
//	@Param			track	body	models.UpdateTrackRequest	true	"Updated track data"
//	@Success		200
//	@Failure		400
//	@Failure		500
//	@Router			/track [put]
func (h *TrackHandler) UpdateTrack(c *gin.Context) {
	var req models.UpdateTrackRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	track := models.Track{
		Id:          req.Id,
		Name:        req.Name,
		Author:      req.Author,
		ReleaseDate: req.ReleaseDate,
		Genres:      req.Genres,
	}

	err := h.trackRepository.Update(&track)
	if err != nil {
		panic(err)
	}

	h.log.Info("Updated a track")
}
