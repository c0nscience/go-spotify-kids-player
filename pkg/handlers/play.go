package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"go-spotify-kids-player/pkg/ha"
	"go-spotify-kids-player/pkg/playlist"
	"go-spotify-kids-player/pkg/store"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
)

type roomSelectionForm struct {
	Rooms []string `form:"rooms[]"`
}

var (
	availableRooms = []string{
		"playroom",
		"kitchen",
		"bathroom",
		"living_room",
	}
	lastPlayedRooms []string
)

func Play(s store.Store) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		objectId, _ := primitive.ObjectIDFromHex(id)

		var form roomSelectionForm
		err := c.ShouldBind(&form)
		if err != nil {
			_ = c.Error(err)
			return
		}
		log.Info().Msgf("selected rooms: %v", form.Rooms)

		if len(form.Rooms) == 0 {
			c.Status(http.StatusNoContent)
			return
		}

		var pl playlist.Playlist
		err = s.FindOne(c, bson.M{"_id": objectId}, &pl)
		if err != nil {
			_ = c.Error(err)
			return
		}

		if len(lastPlayedRooms) > 0 {
			err := ha.Stop(lastPlayedRooms)
			if err != nil {
				_ = c.Error(err)
				return
			}

			err = ha.Unjoin(lastPlayedRooms)
			if err != nil {
				_ = c.Error(err)
				return
			}
		}

		if len(form.Rooms) > 1 {
			err := ha.Join(form.Rooms)
			if err != nil {
				_ = c.Error(err)
				return
			}
		}

		err = ha.Play(pl.Url, form.Rooms)
		if err != nil {
			_ = c.Error(err)
			return
		}

		lastPlayedRooms = form.Rooms

		pl.PlayCount = pl.PlayCount + 1
		pl.Playing = true

		c.Status(http.StatusNoContent)

		var playing []playlist.Playlist
		err = s.Find(c, bson.M{"playing": true}, nil, &playing)
		if err != nil {
			_ = c.Error(err)
			return
		}

		//todo not there yet - we need to avoid that it can be started every where
		// before starting we should stop playing on the other rooms
		// then start the playlist on the newly selected room

		for _, p := range playing {
			p.Playing = false
			_ = s.Save(c, &p)
		}

		err = s.Save(c, &pl)
		if err != nil {
			_ = c.Error(err)
			return
		}

		sendUpdateEvent("")
	}
}
