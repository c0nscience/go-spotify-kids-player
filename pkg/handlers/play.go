package handlers

import (
	"github.com/gin-gonic/gin"
	"go-spotify-kids-player/pkg/ha"
	"go-spotify-kids-player/pkg/playlist"
	"net/http"
)

func Play(c *gin.Context) {
	id := c.Param("id")

	p := playlist.GetById(id)

	err := ha.Play(p.Url)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.Status(http.StatusNoContent)
}
