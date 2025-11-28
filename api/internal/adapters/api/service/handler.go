package serviceapi

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"

	"melodiapp/database"
	serviceports "melodiapp/internal/ports/service"
	"melodiapp/models"
	"melodiapp/shared"
)

type ServiceHandlers struct {
	service serviceports.ServiceService
}

func NewServiceHandlers(s serviceports.ServiceService) *ServiceHandlers {
	return &ServiceHandlers{service: s}
}

func getCurrentUser(c *gin.Context) (*models.User, error) {
	tokenStr := shared.GetTokenFromRequest(c)
	if tokenStr == "" {
		return nil, fmt.Errorf("invalid token")
	}

	token, err := jwt.ParseWithClaims(tokenStr, &shared.Payload{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("invalid token")
		}
		return []byte(os.Getenv("JWT_SECRET_KEY")), nil
	})
	if err != nil {
		return nil, err
	}

	claims, _ := token.Claims.(*shared.Payload)
	session, exists := shared.Sessions[claims.Session]
	if !exists || session.ExpiryTime.Before(time.Now()) {
		return nil, fmt.Errorf("You don't have permission")
	}

	var user models.User
	if err := database.DBConn.First(&user, session.Uid).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

// --- FUNCIÓN AUXILIAR PARA ENRIQUECER SERVICIOS ---
// Esta función hace el trabajo pesado de buscar canciones, usuarios Y OUTFITS
func getServiceDetails(service *models.Service) (gin.H, error) {

	// 1. Obtener canciones asociadas
	var serviceSongs []models.ServiceSong
	if err := database.DBConn.Where("service_id = ?", service.ID).Find(&serviceSongs).Error; err != nil {
		return nil, err
	}

	songIDs := make([]uint, 0, len(serviceSongs))
	for _, ss := range serviceSongs {
		songIDs = append(songIDs, ss.SongID)
	}

	var songs []models.Song
	if len(songIDs) > 0 {
		if err := database.DBConn.Where("id IN ?", songIDs).Find(&songs).Error; err != nil {
			return nil, err
		}
	} else {
		songs = []models.Song{}
	}

	// 2. Obtener usuarios asociados con su status
	var serviceUsers []models.ServiceUser
	if err := database.DBConn.Where("service_id = ?", service.ID).Find(&serviceUsers).Error; err != nil {
		return nil, err
	}

	userIDs := make([]uint, 0, len(serviceUsers))
	statusByUser := make(map[uint]string)
	for _, su := range serviceUsers {
		userIDs = append(userIDs, su.UserID)
		statusByUser[su.UserID] = su.Status
	}

	var users []models.User
	if len(userIDs) > 0 {
		if err := database.DBConn.Where("id IN ?", userIDs).Find(&users).Error; err != nil {
			return nil, err
		}
	}

	// Estructura para el detalle del usuario en la respuesta
	type ServiceUserDetail struct {
		ID            uint   `json:"id"`
		Username      string `json:"username"`
		Email         string `json:"email"`
		Role          string `json:"role"`
		SecondaryRole string `json:"secondary_role"`
		Status        string `json:"status"`
	}

	detailedUsers := make([]ServiceUserDetail, 0, len(users))
	for _, u := range users {
		detailedUsers = append(detailedUsers, ServiceUserDetail{
			ID:            u.ID,
			Username:      u.Username,
			Email:         u.Email,
			Role:          u.Role,
			SecondaryRole: u.SecondaryRole,
			Status:        statusByUser[u.ID],
		})
	}

	// --- 3. NUEVO: Obtener OUTFITS (Paleta de Colores) asociados ---
	var serviceOutfits []models.ServiceOutfit
	if err := database.DBConn.Where("service_id = ?", service.ID).Find(&serviceOutfits).Error; err != nil {
		return nil, err
	}
	// Inicializar como array vacío si es nil para evitar problemas en frontend
	if serviceOutfits == nil {
		serviceOutfits = []models.ServiceOutfit{}
	}

	// Construir el objeto de respuesta unificado
	return gin.H{
		"id":         service.ID,
		"name":       service.Name,
		"start_time": service.StartTime,
		"end_time":   service.EndTime,
		"created_by": service.CreatedBy,
		"created_at": service.CreatedAt,
		"updated_at": service.UpdatedAt,
		"songs":      songs,
		"users":      detailedUsers,
		"outfits":    serviceOutfits, // <--- CAMPO AGREGADO
	}, nil
}

// --- HANDLERS ---

func (h *ServiceHandlers) GetAll(c *gin.Context) {
	if _, err := getCurrentUser(c); err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	services, err := h.service.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var fullServices []gin.H
	for _, s := range services {
		details, err := getServiceDetails(&s)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error loading details: " + err.Error()})
			return
		}
		fullServices = append(fullServices, details)
	}

	if fullServices == nil {
		fullServices = []gin.H{}
	}

	c.JSON(http.StatusOK, fullServices)
}

func (h *ServiceHandlers) GetByID(c *gin.Context) {
	if _, err := getCurrentUser(c); err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	id := c.Param("id")
	service, err := h.service.GetByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if service == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Service not found"})
		return
	}

	response, err := getServiceDetails(service)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	nestedResponse := gin.H{
		"service": service,
		"songs":   response["songs"],
		"users":   response["users"],
		"outfits": response["outfits"], // <--- CAMPO AGREGADO
	}

	c.JSON(http.StatusOK, nestedResponse)
}

func (h *ServiceHandlers) Create(c *gin.Context) {
	user, err := getCurrentUser(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	if user.Role != "admin" {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "You don't have permission"})
		return
	}

	var input models.Service
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data"})
		return
	}

	input.CreatedBy = user.ID

	created, err := h.service.Create(&input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, created)
}

func (h *ServiceHandlers) Update(c *gin.Context) {
	user, err := getCurrentUser(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	if user.Role != "admin" {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "You don't have permission"})
		return
	}

	id := c.Param("id")
	var input models.Service
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data"})
		return
	}

	updated, err := h.service.Update(id, &input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if updated == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Service not found"})
		return
	}

	c.JSON(http.StatusOK, updated)
}

func (h *ServiceHandlers) Delete(c *gin.Context) {
	user, err := getCurrentUser(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	if user.Role != "admin" {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "You don't have permission"})
		return
	}

	id := c.Param("id")
	if err := h.service.Delete(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusNoContent, gin.H{"id": id})
}
