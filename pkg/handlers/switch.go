package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"net/http"
	"time"
)

func Switch() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Query("id")
		room := c.Query("room")
		log.Info().Msgf("go switch %s to %s", id, room)
		time.Sleep(1 * time.Second)
		c.Status(http.StatusNoContent)

		/**
		todo
			group current room with target room
			wait for target room to start playing
			wait another grace-period
			remove previous room from group
		*/

	}
}
