package handlers

import (
	"github.com/gin-gonic/gin"
	"go-spotify-kids-player/pkg/playlist"
	"net/http"
)

func List(c *gin.Context) {
	playlists := playlist.GetAll()

	c.HTML(http.StatusOK, "list.gohtml", gin.H{
		"Playlists": playlists,
	})
}
