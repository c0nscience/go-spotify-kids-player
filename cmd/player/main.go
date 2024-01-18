package main

import (
	"context"
	"embed"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	spotifyapi "github.com/zmb3/spotify/v2"
	"github.com/zmb3/spotify/v2/auth"
	"go-spotify-kids-player/pkg/handlers"
	"go-spotify-kids-player/pkg/playlist"
	"go-spotify-kids-player/pkg/store"
	"golang.org/x/oauth2/clientcredentials"
	"html/template"
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

	dbUri := os.Getenv("DB_URI")
	dbName := os.Getenv("DB_NAME")
	playlistClient, err := store.New(dbUri, dbName, playlist.Collection)
	if err != nil {
		log.Panic().Err(err).Msg("Could not create activity store")
	}

	r := gin.Default()
	r.SetFuncMap(template.FuncMap{
		"join": strings.Join,
	})
	r.LoadHTMLGlob("templates/*.gohtml")

	r.StaticFS("/public", http.FS(f))

	r.GET("/", handlers.List(playlistClient))
	r.GET("/:id/play", handlers.Play(playlistClient))
	r.GET("/edit", handlers.Edit(playlistClient))
	r.POST("/add", handlers.Add(syCli, playlistClient))
	r.DELETE("/:id/delete", handlers.Delete(playlistClient))

	srv := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatal().Err(err).Msgf("listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Log().Msg("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal().Err(err).Msgf("Server Shutdown failed")
	}
	if err := playlistClient.Disconnect(ctx); err != nil {
		log.Fatal().Err(err).Msg("Disconnect from mongodb failed")
	}

	select {
	case <-ctx.Done():
		log.Log().Msg("timeout of 5 seconds.")
	}
	log.Log().Msg("Server exiting")

}
