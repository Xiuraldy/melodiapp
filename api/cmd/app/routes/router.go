package routes

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	authroutes "melodiapp/cmd/app/routes/auth"
	serviceroutes "melodiapp/cmd/app/routes/service"
	songroutes "melodiapp/cmd/app/routes/song"
	userroutes "melodiapp/cmd/app/routes/user"
	"melodiapp/database"
	"melodiapp/shared"
)

func NewRouter() *gin.Engine {
	r := gin.Default()
	r.Use(shared.Cors())

	r.Static("/files", "./public")

	userroutes.AddUserRoutes(r)
	authroutes.AddAuthRoutes(r)
	serviceroutes.AddServiceRoutes(r)
	songroutes.AddSongRoutes(r)

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
