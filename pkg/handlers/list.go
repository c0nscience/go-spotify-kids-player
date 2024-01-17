package handlers

import (
	"github.com/gin-gonic/gin"
	"go-spotify-kids-player/pkg/models"
	"go-spotify-kids-player/pkg/playlist"
	"go-spotify-kids-player/pkg/spotify"
	"go-spotify-kids-player/web/template"
	"net/http"
	"os"
)

func List(c *gin.Context) {
	accessToken, err := spotify.GetAccessToken(os.Getenv("SPOTIFY_CLIENT_ID"), os.Getenv("SPOTIFY_CLIENT_SECRET"))
	if err != nil {
		_ = c.Error(err)
		return
	}

	paylists := playlist.GetAll()

	playlistViewModels := []models.PlaylistViewModel{}

	for _, p := range paylists {
		id := models.GetId(p)

		img, err := spotify.GetCoverImage(accessToken, id)
		if err == nil {
			playlistViewModels = append(playlistViewModels, models.PlaylistViewModel{
				Img: img,
				Id:  p.Id,
			})
		}
	}

	c.HTML(http.StatusOK, "", template.List(playlistViewModels))
}
