package handlers

import (
	"github.com/gin-gonic/gin"
	spotifyapi "github.com/zmb3/spotify/v2"
	"go-spotify-kids-player/pkg/playlist"
	"go-spotify-kids-player/pkg/spotify"
	"net/http"
)

func Edit(c *gin.Context) {
	playlists := playlist.GetAll()

	c.HTML(http.StatusOK, "edit.gohtml", gin.H{
		"Playlists": playlists,
	})
}

func Add(cli *spotifyapi.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		u := c.PostForm("url")
		spotifyId := spotify.GetIdFrom(u)

		album, err := spotify.GetAlbum(c, cli, spotifyId)
		if err != nil {
			_ = c.Error(err)
			return
		}

		imgUrl := ""
		if len(album.Images) > 1 {
			imgUrl = album.Images[0].URL
		}

		var artistNames []string
		for _, artist := range album.Artists {
			artistNames = append(artistNames, artist.Name)
		}

		storedPlaylist := playlist.Store(playlist.Playlist{
			Url:       u,
			Name:      album.Name,
			Img:       imgUrl,
			SpotifyID: spotifyId,
			Artists:   artistNames,
		})

		c.HTML(http.StatusOK, "edit-list-entry", gin.H{
			"ID":      storedPlaylist.ID,
			"Img":     storedPlaylist.Img,
			"Name":    storedPlaylist.Name,
			"Artists": storedPlaylist.Artists,
		})
	}
}

func Delete(c *gin.Context) {
	id := c.Param("id")

	playlist.DeleteByID(id)
	c.Status(http.StatusOK)
}
