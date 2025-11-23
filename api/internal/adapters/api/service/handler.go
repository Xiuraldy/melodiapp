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
	c.JSON(http.StatusOK, services)
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

	// Obtener canciones asociadas al servicio
	var serviceSongs []models.ServiceSong
	if err := database.DBConn.Where("service_id = ?", service.ID).Find(&serviceSongs).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	songIDs := make([]uint, 0, len(serviceSongs))
	for _, ss := range serviceSongs {
		songIDs = append(songIDs, ss.SongID)
	}

	var songs []models.Song
	if len(songIDs) > 0 {
		if err := database.DBConn.Where("id IN ?", songIDs).Find(&songs).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	// Obtener usuarios asociados al servicio con su status
	var serviceUsers []models.ServiceUser
	if err := database.DBConn.Where("service_id = ?", service.ID).Find(&serviceUsers).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	userIDs := make([]uint, 0, len(serviceUsers))
	for _, su := range serviceUsers {
		userIDs = append(userIDs, su.UserID)
	}

	var users []models.User
	if len(userIDs) > 0 {
		if err := database.DBConn.Where("id IN ?", userIDs).Find(&users).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	type ServiceUserDetail struct {
		ID       uint   `json:"id"`
		Username string `json:"username"`
		Email    string `json:"email"`
		Role     string `json:"role"`
		Status   string `json:"status"`
	}

	statusByUser := make(map[uint]string)
	for _, su := range serviceUsers {
		statusByUser[su.UserID] = su.Status
	}

	detailedUsers := make([]ServiceUserDetail, 0, len(users))
	for _, u := range users {
		detailedUsers = append(detailedUsers, ServiceUserDetail{
			ID:       u.ID,
			Username: u.Username,
			Email:    u.Email,
			Role:     u.Role,
			Status:   statusByUser[u.ID],
		})
	}

	response := gin.H{
		"service": service,
		"songs":   songs,
		"users":   detailedUsers,
	}

	c.JSON(http.StatusOK, response)
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

	// Forzar created_by al usuario autenticado, ignorando lo que venga del body
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
