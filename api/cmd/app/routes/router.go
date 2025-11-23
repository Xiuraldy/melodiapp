package routes

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	authroutes "melodiapp/cmd/app/routes/auth"
	userroutes "melodiapp/cmd/app/routes/user"
	"melodiapp/database"
	"melodiapp/shared"
)

func NewRouter() *gin.Engine {
	r := gin.Default()
	r.Use(shared.Cors())

	userroutes.AddUserRoutes(r)
	authroutes.AddAuthRoutes(r)

	r.GET("/", func(c *gin.Context) {
		tx := database.DBConn.Exec("SELECT 1")
		if tx.Error != nil {
			log.Printf("Error al ejecutar consulta: %v", tx.Error)
			c.JSON(http.StatusInternalServerError, gin.H{"Success": false})
			return
		}
		c.JSON(http.StatusOK, gin.H{"Success": true})
	})

	return r
}
