package userapi

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"

	"melodiapp/database"
	userports "melodiapp/internal/ports/user"
	"melodiapp/models"
	"melodiapp/shared"
)

type UserHandlers struct {
	service userports.UserService
}

func NewUserHandlers(service userports.UserService) *UserHandlers {
	return &UserHandlers{service: service}
}

func (h *UserHandlers) GetAllUsers(c *gin.Context) {
	tokenStr := shared.GetTokenFromRequest(c)
	token, err := jwt.ParseWithClaims(tokenStr, &shared.Payload{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("invalid token")
		}
		return []byte(os.Getenv("JWT_SECRET_KEY")), nil
	})

	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	claims, _ := token.Claims.(*shared.Payload)
	userData, exists := shared.Sessions[claims.Session]

	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "You don't have permission"})
		return
	}

	_ = userData // authentication already validated, we don't use userData directly here

	users, err := h.service.GetAllUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, users)
}

func (h *UserHandlers) GetMe(c *gin.Context) {
	tokenStr := shared.GetTokenFromRequest(c)
	token, err := jwt.ParseWithClaims(tokenStr, &shared.Payload{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("invalid token")
		}
		return []byte(os.Getenv("JWT_SECRET_KEY")), nil
	})

	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	claims, _ := token.Claims.(*shared.Payload)
	userData, exists := shared.Sessions[claims.Session]

	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "You don't have permission"})
		return
	}

	if userData.ExpiryTime.Before(time.Now()) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "You don't have permission"})
		return
	}

	user, err := h.service.GetUserByUintID(userData.Uid)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "You don't have permission"})
		return
	}
	if user == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}

func (h *UserHandlers) GetUserById(c *gin.Context) {
	tokenStr := shared.GetTokenFromRequest(c)
	token, err := jwt.ParseWithClaims(tokenStr, &shared.Payload{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("invalid token")
		}
		return []byte(os.Getenv("JWT_SECRET_KEY")), nil
	})

	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	claims, _ := token.Claims.(*shared.Payload)
	userData, exists := shared.Sessions[claims.Session]

	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "You don't have permission"})
		return
	}

	if userData.ExpiryTime.Before(time.Now()) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "You don't have permission"})
		return
	}

	// ensure the session user still exists
	if _, err := h.service.GetUserByUintID(userData.Uid); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "You don't have permission"})
		return
	}

	id := c.Param("id")
	user, err := h.service.GetUserByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if user == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}

func (h *UserHandlers) CreateUser(c *gin.Context) {
	// Crear usuario a partir de campos de formulario (multipart/form-data)
	username := c.PostForm("username")
	lastname := c.PostForm("lastname")
	email := c.PostForm("email")
	celphone := c.PostForm("celphone")
	role := c.PostForm("role")
	secondaryRole := c.PostForm("secondary_role")
	password := c.PostForm("password")

	if username == "" || email == "" || password == "" || role == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "username, email, password y role son obligatorios"})
		return
	}

	user := models.User{
		Username:      username,
		Lastname:      lastname,
		Email:         email,
		Celphone:      celphone,
		Role:          role,
		SecondaryRole: secondaryRole,
		Password:      password,
	}

	// Manejar archivo de foto opcional
	file, err := c.FormFile("file")
	if err == nil {
		uploadDir := "public/profiles"
		if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
			os.MkdirAll(uploadDir, 0755)
		}

		filename := fmt.Sprintf("profile_%d_%s", time.Now().UnixNano(), file.Filename)
		filepath := fmt.Sprintf("%s/%s", uploadDir, filename)

		if err := c.SaveUploadedFile(file, filepath); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save image"})
			return
		}

		user.ProfilePictureUrl = fmt.Sprintf("/files/profiles/%s", filename)
	}

	created, err := h.service.CreateUser(&user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, created)
}

func (h *UserHandlers) EditUser(c *gin.Context) {
	// 1. VALIDACIÓN DE TOKEN (Tu lógica actual)
	tokenStr := shared.GetTokenFromRequest(c)
	token, err := jwt.ParseWithClaims(tokenStr, &shared.Payload{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("invalid token")
		}
		return []byte(os.Getenv("JWT_SECRET_KEY")), nil
	})

	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	claims, _ := token.Claims.(*shared.Payload)
	_, exists := shared.Sessions[claims.Session]
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "You don't have permission"})
		return
	}

	// 2. BUSCAR EL USUARIO A EDITAR
	id := c.Param("id")
	var user models.User
	tx := database.DBConn.First(&user, id)
	if tx.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// 3. OBTENER DATOS DEL FORMULARIO (MULTIPART/FORM-DATA)
	// Nota: Ya no usamos BindJSON porque vienen archivos

	// --- Campos de Texto ---
	username := c.PostForm("username")
	lastname := c.PostForm("lastname")
	email := c.PostForm("email")
	celphone := c.PostForm("celphone")
	role := c.PostForm("role")
	secondaryRole := c.PostForm("secondary_role")
	password := c.PostForm("password")

	// Actualizar campos si vienen datos
	if username != "" {
		user.Username = username
	}
	if lastname != "" {
		user.Lastname = lastname
	}
	if email != "" {
		user.Email = email
	}
	if celphone != "" {
		user.Celphone = celphone
	}
	if role != "" {
		user.Role = role
	}
	if secondaryRole != "" {
		user.SecondaryRole = secondaryRole
	}

	// --- Contraseña (Solo si se envía) ---
	// Asignamos la nueva contraseña en texto plano y dejamos que los hooks de GORM
	// (BeforeUpdate) se encarguen de hashearla si realmente cambió.
	if password != "" {
		user.Password = password
	}

	// --- 4. MANEJO DEL ARCHIVO (FOTO) ---
	// El frontend envía el archivo en el campo "file"
	file, err := c.FormFile("file")
	if err == nil {
		// Si hay archivo, lo guardamos

		// Asegurar que el directorio existe
		uploadDir := "public/profiles"
		if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
			os.MkdirAll(uploadDir, 0755)
		}

		// Crear nombre único: profile_{id}_{originalName}
		filename := fmt.Sprintf("profile_%s_%s", id, file.Filename)
		filepath := fmt.Sprintf("%s/%s", uploadDir, filename)

		// Guardar en disco
		if err := c.SaveUploadedFile(file, filepath); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save image"})
			return
		}

		// Actualizar la URL en la base de datos
		// La ruta web será /files/profiles/... (dependiendo de cómo sirvas los estáticos)
		user.ProfilePictureUrl = fmt.Sprintf("/files/profiles/%s", filename)
	}

	// 5. GUARDAR CAMBIOS EN BD
	if err := database.DBConn.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "User updated successfully",
		"user":    user,
	})
}

func (h *UserHandlers) DeleteUser(c *gin.Context) {
	id := c.Param("id")
	if err := h.service.DeleteUser(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusNoContent, gin.H{"id": id})
}
