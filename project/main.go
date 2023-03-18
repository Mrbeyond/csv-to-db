package main

import (
	"context"
	"csvapi-test/model"
	"csvapi-test/router"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		panic("ENV not loaded")
	}

	// Establish database connection,
	model.DbConfig()

	// Go gin Default engine visit https://github.com/gin-gonic/gin
	app := gin.Default()

	app.GET("", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, "Up and running")
	})

	// Add app controller routers
	router.AppInstance(app)

	app.NoRoute(func(ctx *gin.Context) {
		notFoundResponse := gin.H{
			"error":   true,
			"status":  "failed",
			"message": "Page not found",
		}
		ctx.JSON(http.StatusNotFound, notFoundResponse)
	})

	server := &http.Server{
		Addr:           fmt.Sprintf(":%s", os.Getenv("PORT")),
		Handler:        app,
		ReadTimeout:    1 * time.Minute, // Default is 10s, increaased incase of large files
		WriteTimeout:   3 * time.Minute, // Default is 10s, increaased incase of longer operation
		MaxHeaderBytes: 2 << 20,         //2MB
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("listen:", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with a timeout context.
	// Visit https://gin-gonic.com/docs/examples/graceful-restart-or-stop/
	quit := make(chan os.Signal)
	// kill (no param) default send syscanll.SIGTERM
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatal("Server shutdown", err)
	}

	// Catching ctx.Done().
	select {
	case <-ctx.Done():
		log.Println("Graceful shutdown completed")
	}
}
