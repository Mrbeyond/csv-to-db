package router

import (
	"csvapi-test/controller"
	"csvapi-test/middleware"

	"github.com/gin-gonic/gin"
)

// App server engine instance with registered routes
func AppInstance(app *gin.Engine) {

	app.Use(middleware.CORSMiddleware())

	app.Use(middleware.TimeoutMiddleware())

	app.POST("/data", controller.Create)
	app.GET("/data", controller.Fetch)
}
