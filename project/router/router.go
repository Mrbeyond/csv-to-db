package router

import (
	"csvapi-test/controller"
	"csvapi-test/middleware"

	"github.com/gin-gonic/gin"
)

// App server engine instance with registered routes
func AppInstance() *gin.Engine {
	// Go gin Default engine visit https://github.com/gin-gonic/gin
	app := gin.Default()

	app.Use(middleware.CORSMiddleware())
	app.Use(middleware.TimeoutMiddleware())

	app.POST("/data", controller.Create)
	app.GET("/data", controller.Fetch)

	return app
}
