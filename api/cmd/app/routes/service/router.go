package service

import (
	"github.com/gin-gonic/gin"

	serviceapi "melodiapp/internal/adapters/api/service"
	serviceoutfitapi "melodiapp/internal/adapters/api/serviceoutfit"
	servicesongapi "melodiapp/internal/adapters/api/servicesong"
	serviceuserapi "melodiapp/internal/adapters/api/serviceuser"

	dbadapter "melodiapp/internal/adapters/database/service"
	dbserviceoutfit "melodiapp/internal/adapters/database/serviceoutfit"
	dbservicesong "melodiapp/internal/adapters/database/servicesong"
	dbserviceuser "melodiapp/internal/adapters/database/serviceuser"

	coreservice "melodiapp/internal/core/service"
	coreservicesong "melodiapp/internal/core/servicesong"
	coreserviceuser "melodiapp/internal/core/serviceuser"

	// 3. Import del Core
	coreserviceoutfit "melodiapp/internal/core/serviceoutfit"
)

func AddServiceRoutes(r *gin.Engine) {
	group := r.Group("/services")

	serviceRepo := dbadapter.NewGormServiceRepository()
	serviceUsecase := coreservice.NewServiceUsecase(serviceRepo)
	serviceHandlers := serviceapi.NewServiceHandlers(serviceUsecase)

	serviceUserRepo := dbserviceuser.NewGormServiceUserRepository()
	serviceUserUsecase := coreserviceuser.NewService(serviceUserRepo)
	serviceUserHandlers := serviceuserapi.NewServiceUserHandlers(serviceUserUsecase)

	serviceSongRepo := dbservicesong.NewGormServiceSongRepository()
	serviceSongUsecase := coreservicesong.NewService(serviceSongRepo)
	serviceSongHandlers := servicesongapi.NewServiceSongHandlers(serviceSongUsecase)

	serviceOutfitRepo := dbserviceoutfit.NewGormServiceOutfitRepository()
	serviceOutfitUsecase := coreserviceoutfit.NewService(serviceOutfitRepo)
	serviceOutfitHandlers := serviceoutfitapi.NewServiceOutfitHandlers(serviceOutfitUsecase)

	group.GET("", serviceHandlers.GetAll)
	group.GET(":id", serviceHandlers.GetByID)
	group.POST("", serviceHandlers.Create)
	group.PUT(":id", serviceHandlers.Update)
	group.DELETE(":id", serviceHandlers.Delete)

	group.POST(":id/users", serviceUserHandlers.AssignUsers)
	group.GET(":id/users", serviceUserHandlers.ListByService)
	group.PATCH(":id/users/:userId/status", serviceUserHandlers.ChangeStatus)

	group.POST(":id/songs", serviceSongHandlers.AssignSongs)
	group.GET(":id/songs", serviceSongHandlers.ListByService)
	group.DELETE(":id/songs/:songId", serviceSongHandlers.Remove)

	group.POST(":id/outfits", serviceOutfitHandlers.AssignOutfits)
	group.GET(":id/outfits", serviceOutfitHandlers.ListByService)
	group.DELETE(":id/outfits/:outfitId", serviceOutfitHandlers.Remove)
}
