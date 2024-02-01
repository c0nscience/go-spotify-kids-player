package handlers

import (
	"github.com/gin-gonic/gin"
	"go-spotify-kids-player/pkg/player"
	"go-spotify-kids-player/pkg/playlist"
	"go-spotify-kids-player/pkg/store"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
)

type PlaylistListViewModel struct {
	ID      string
	Img     string
	Playing bool
}

type PlayerViewModel struct {
	ID      string
	Img     string
	Playing bool
	Rooms   []string
}

func ListView(s store.Store) gin.HandlerFunc {
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
				ID:      p.ID.Hex(),
				Img:     p.Img,
				Playing: p.Playing,
			})
		}

		state := player.State()
		c.HTML(http.StatusOK, "list.gohtml", gin.H{
			"Playlists": viewModels,
			"Player": PlayerViewModel{
				ID:      state.ID,
				Playing: state.Playing,
				Img:     state.Img,
				Rooms:   availableRooms,
			},
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
				ID:      p.ID.Hex(),
				Img:     p.Img,
				Playing: p.Playing,
			})
		}

		c.HTML(http.StatusOK, "playlist-list", gin.H{
			"Playlists": viewModels,
		})
	}
}

func RoomSelectionModal(s store.Store) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		objectId, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			_ = c.Error(err)
			return
		}

		var pl playlist.Playlist
		err = s.FindOne(c, bson.M{"_id": objectId}, &pl)
		if err != nil {
			_ = c.Error(err)
			return
		}

		c.HTML(http.StatusOK, "room-selection-modal.gohtml", gin.H{
			"ID":    id,
			"Img":   pl.Img,
			"Rooms": availableRooms,
		})
	}
}
