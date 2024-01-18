package handlers

import (
	"github.com/gin-gonic/gin"
	"go-spotify-kids-player/pkg/playlist"
	"go-spotify-kids-player/pkg/store"
	"go.mongodb.org/mongo-driver/bson"
	"net/http"
)

func List(s store.Store) gin.HandlerFunc {
	return func(c *gin.Context) {
		var playlists []playlist.Playlist

		err := s.Find(c, bson.D{}, nil, &playlists)
		if err != nil {
			_ = c.Error(err)
			return
		}

		c.HTML(http.StatusOK, "list.gohtml", gin.H{
			"Playlists": playlists,
		})
	}
}
