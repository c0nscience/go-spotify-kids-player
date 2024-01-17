package handlers

import (
	"github.com/gin-gonic/gin"
	"go-spotify-kids-player/pkg/playlist"
	"log"
	"net/http"
)

func Edit(c *gin.Context) {
	playlists := playlist.GetAll()

	c.HTML(http.StatusOK, "edit.gohtml", gin.H{
		"Playlists": playlists,
	})
}

func Add(c *gin.Context) {
	u := c.PostForm("url")
	log.Printf("add %s", u)

	c.Status(http.StatusNoContent)
}
