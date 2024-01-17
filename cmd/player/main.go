package main

import (
	"context"
	"embed"
	"errors"
	"github.com/gin-gonic/gin"
	"go-spotify-kids-player/pkg/handlers"
	"go-spotify-kids-player/pkg/renderer"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

//go:embed assets/*
var f embed.FS

func main() {
	r := gin.Default()
	r.HTMLRender = &renderer.TemplRender{}

	r.StaticFS("/public", http.FS(f))

	r.GET("/", handlers.List)
	r.GET("/:id/play", handlers.Play)

	srv := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	go func() {
		// service connections
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	select {
	case <-ctx.Done():
		log.Println("timeout of 5 seconds.")
	}
	log.Println("Server exiting")

}
