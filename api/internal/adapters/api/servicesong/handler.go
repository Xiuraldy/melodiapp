package servicesongapi

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	servicesongports "melodiapp/internal/ports/servicesong"
)

type ServiceSongHandlers struct {
	service servicesongports.ServiceSongService
}

type assignSongsRequest struct {
	SongIDs []uint `json:"song_ids"`
}

func NewServiceSongHandlers(s servicesongports.ServiceSongService) *ServiceSongHandlers {
	return &ServiceSongHandlers{service: s}
}

func (h *ServiceSongHandlers) AssignSongs(c *gin.Context) {
	serviceIDParam := c.Param("id")
	serviceID64, err := strconv.ParseUint(serviceIDParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid service id"})
		return
	}

	var req assignSongsRequest
	if err := c.ShouldBindJSON(&req); err != nil || len(req.SongIDs) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data"})
		return
	}

	if err := h.service.AssignSongs(uint(serviceID64), req.SongIDs); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"service_id": serviceIDParam, "song_ids": req.SongIDs})
}

func (h *ServiceSongHandlers) ListByService(c *gin.Context) {
	serviceIDParam := c.Param("id")
	serviceID64, err := strconv.ParseUint(serviceIDParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid service id"})
		return
	}

	items, err := h.service.ListByService(uint(serviceID64))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, items)
}

func (h *ServiceSongHandlers) Remove(c *gin.Context) {
	serviceIDParam := c.Param("id")
	songIDParam := c.Param("songId")

	serviceID64, err := strconv.ParseUint(serviceIDParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid service id"})
		return
	}
	songID64, err := strconv.ParseUint(songIDParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid song id"})
		return
	}

	if err := h.service.Remove(uint(serviceID64), uint(songID64)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusNoContent, gin.H{"service_id": serviceIDParam, "song_id": songIDParam})
}
