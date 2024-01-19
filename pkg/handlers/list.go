package handlers

import (
	"github.com/gin-gonic/gin"
	"go-spotify-kids-player/pkg/playlist"
	"go-spotify-kids-player/pkg/store"
	"go.mongodb.org/mongo-driver/bson"
	"net/http"
)

type PlaylistListViewModel struct {
	ID  string
	Img string
}

func List(s store.Store) gin.HandlerFunc {
	return func(c *gin.Context) {
		var playlists []playlist.Playlist

		err := s.Find(c, bson.D{}, nil, &playlists)
		if err != nil {
			_ = c.Error(err)
			return
		}

		var viewModels []PlaylistListViewModel
		for _, p := range playlists {
			viewModels = append(viewModels, PlaylistListViewModel{
				ID:  p.ID.Hex(),
				Img: p.Img,
			})
		}

		c.HTML(http.StatusOK, "list.gohtml", gin.H{
			"Playlists": viewModels,
		})
	}
}

func PartialList(s store.Store) gin.HandlerFunc {
	return func(c *gin.Context) {
		var playlists []playlist.Playlist

		err := s.Find(c, bson.D{}, nil, &playlists)
		if err != nil {
			_ = c.Error(err)
			return
		}

		var viewModels []PlaylistListViewModel
		for _, p := range playlists {
			viewModels = append(viewModels, PlaylistListViewModel{
				ID:  p.ID.Hex(),
				Img: p.Img,
			})
		}

		c.HTML(http.StatusOK, "playlist-list", gin.H{
			"Playlists": viewModels,
		})
	}
}
