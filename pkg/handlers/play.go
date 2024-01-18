package handlers

import (
	"github.com/gin-gonic/gin"
	"go-spotify-kids-player/pkg/ha"
	"go-spotify-kids-player/pkg/playlist"
	"go-spotify-kids-player/pkg/store"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
)

func Play(s store.Store) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		objectId, _ := primitive.ObjectIDFromHex(id)

		var pl playlist.Playlist
		err := s.FindOne(c, bson.M{"_id": objectId}, &pl)
		if err != nil {
			_ = c.Error(err)
			return
		}

		err = ha.Play(pl.Url)
		if err != nil {
			_ = c.Error(err)
			return
		}

		c.Status(http.StatusNoContent)
	}
}
