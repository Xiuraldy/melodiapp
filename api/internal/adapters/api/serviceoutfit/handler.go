package serviceoutfitapi

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	// Asumo que crearás este puerto/interface similar al de songs
	serviceoutfitports "melodiapp/internal/ports/serviceoutfit"
)

type ServiceOutfitHandlers struct {
	service serviceoutfitports.ServiceOutfitService
}

// Estructura para recibir los IDs de los outfits en el body del request
type assignOutfitsRequest struct {
	OutfitIDs []uint `json:"outfit_ids"`
}

func NewServiceOutfitHandlers(s serviceoutfitports.ServiceOutfitService) *ServiceOutfitHandlers {
	return &ServiceOutfitHandlers{service: s}
}

// AssignOutfits vincula uno o más outfits a un servicio
func (h *ServiceOutfitHandlers) AssignOutfits(c *gin.Context) {
	serviceIDParam := c.Param("id")
	serviceID64, err := strconv.ParseUint(serviceIDParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid service id"})
		return
	}

	var req assignOutfitsRequest
	// Valida que el JSON sea correcto y que el array no esté vacío
	if err := c.ShouldBindJSON(&req); err != nil || len(req.OutfitIDs) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data, outfit_ids required"})
		return
	}

	// Llama al servicio para guardar la relación en la tabla service_outfit
	if err := h.service.AssignOutfits(uint(serviceID64), req.OutfitIDs); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"service_id": serviceIDParam, "outfit_ids": req.OutfitIDs})
}

// ListByService obtiene todos los outfits asignados a un servicio
func (h *ServiceOutfitHandlers) ListByService(c *gin.Context) {
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

// Remove elimina la relación entre un servicio y un outfit específico
func (h *ServiceOutfitHandlers) Remove(c *gin.Context) {
	serviceIDParam := c.Param("id")
	outfitIDParam := c.Param("outfitId") // Asegúrate de usar este param en la ruta de Gin

	serviceID64, err := strconv.ParseUint(serviceIDParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid service id"})
		return
	}
	outfitID64, err := strconv.ParseUint(outfitIDParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid outfit id"})
		return
	}

	if err := h.service.Remove(uint(serviceID64), uint(outfitID64)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusNoContent, gin.H{"service_id": serviceIDParam, "outfit_id": outfitIDParam})
}
