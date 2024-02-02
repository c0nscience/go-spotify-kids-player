package handlers

import (
	"github.com/gin-gonic/gin"
	"go-spotify-kids-player/pkg/player"
	"net/http"
)

func Player() gin.HandlerFunc {
	return func(c *gin.Context) {

		state := player.State()
		c.HTML(http.StatusOK, "player", gin.H{
			"Player": PlayerViewModel{
				ID:      state.ID,
				Playing: state.Playing,
				Img:     state.Img,
				Rooms:   availableRooms,
			},
		})
	}
}
