package main

import (
	"context"
	"embed"
	"errors"
	"github.com/gin-gonic/gin"
	spotifyapi "github.com/zmb3/spotify/v2"
	"github.com/zmb3/spotify/v2/auth"
	"go-spotify-kids-player/pkg/handlers"
	"golang.org/x/oauth2/clientcredentials"
	"html/template"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

//go:embed assets/*
var f embed.FS

func main() {
	ctx := context.Background()
	config := &clientcredentials.Config{
		ClientID:     os.Getenv("SPOTIFY_ID"),
		ClientSecret: os.Getenv("SPOTIFY_SECRET"),
		TokenURL:     spotifyauth.TokenURL,
	}
	httpClient := config.Client(ctx)
	syCli := spotifyapi.New(httpClient)

	r := gin.Default()
	r.SetFuncMap(template.FuncMap{
		"join": strings.Join,
	})
	r.LoadHTMLGlob("templates/*.gohtml")

	r.StaticFS("/public", http.FS(f))

	r.GET("/", handlers.List)
	r.GET("/:id/play", handlers.Play)
	r.GET("/edit", handlers.Edit)
	r.POST("/add", handlers.Add(syCli))
	r.DELETE("/:id/delete", handlers.Delete)

	srv := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
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
