package handlers

import (
	"github.com/gin-gonic/gin"
	spotify2 "github.com/zmb3/spotify/v2"
	"go-spotify-kids-player/pkg/playlist"
	"go-spotify-kids-player/pkg/spotify"
	"net/http"
)

type playlistViewModel struct {
	Img string
	Id  string
}

func List(cli *spotify2.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		playlists := playlist.GetAll()

		var playlistViewModels []playlistViewModel

		for _, p := range playlists {
			id := playlist.GetId(p)

			img, err := spotify.GetCoverImage(c, cli, id)
			if err == nil {
				playlistViewModels = append(playlistViewModels, playlistViewModel{
					Img: img,
					Id:  p.Id,
				})
			}
		}

		c.HTML(http.StatusOK, "list.gohtml", gin.H{
			"Playlists": playlistViewModels,
		})
	}
}
