package main

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
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

func init() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339})
}

func main() {
	//dir, _ := os.Getwd()
	//log.Info().Msgf("pwd: %s", dir)
	//return
	ctx := context.Background()
	config := &clientcredentials.Config{
		ClientID:     os.Getenv("SPOTIFY_ID"),
		ClientSecret: os.Getenv("SPOTIFY_SECRET"),
		TokenURL:     spotifyauth.TokenURL,
	}
	httpClient := config.Client(ctx)
	syCli := spotifyapi.New(httpClient, spotifyapi.WithAcceptLanguage("DE"))

	dbUri := os.Getenv("DB_URI")
	dbName := os.Getenv("DB_NAME")
	playlistClient, err := store.New(dbUri, dbName, playlist.Collection)
	if err != nil {
		log.Panic().Err(err).Msg("Could not create activity store")
	}

	r := gin.Default()
	r.Use(handlers.ErrorHandling)

	r.SetFuncMap(template.FuncMap{
		"join": strings.Join,
	})
	r.LoadHTMLGlob("./templates/*.gohtml")

	r.Static("/public", "./public")

	r.GET("/", handlers.ListView(playlistClient))
	r.POST("/:id/play", handlers.Play(playlistClient))
	r.GET("/edit", handlers.EditView(playlistClient))
	r.POST("/add", handlers.Add(syCli, playlistClient))
	r.DELETE("/:id/delete", handlers.Delete(playlistClient))
	r.GET("/sse", handlers.SSE)
	r.GET("/update-list", handlers.PartialList(playlistClient))
	r.GET("/:id/select-room", handlers.RoomSelectionModal())

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
	log.Info().Msg("Shutdown Server ...")

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
		log.Info().Msg("timeout of 5 seconds.")
	}
	log.Info().Msg("Server exiting")

}
