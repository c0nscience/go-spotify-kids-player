package handlers

import (
	"github.com/gin-gonic/gin"
	spotifyapi "github.com/zmb3/spotify/v2"
	"go-spotify-kids-player/pkg/playlist"
	"go-spotify-kids-player/pkg/spotify"
	"go-spotify-kids-player/pkg/sse"
	"go-spotify-kids-player/pkg/store"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"strings"
)

type EditListViewModel struct {
	ID        string
	Img       string
	Name      string
	Artists   []string
	PlayCount int
}

func EditView(s store.Store) gin.HandlerFunc {
	return func(c *gin.Context) {
		var playlists []playlist.Playlist

		err := s.Find(c, bson.D{}, nil, &playlists)
		if err != nil {
			_ = c.Error(err)
			return
		}

		var viewModels []EditListViewModel
		for _, p := range playlists {
			viewModels = append(viewModels, EditListViewModel{
				ID:        p.ID.Hex(),
				Img:       p.Img,
				Name:      p.Name,
				Artists:   p.Artists,
				PlayCount: p.PlayCount,
			})
		}

		c.HTML(http.StatusOK, "edit.gohtml", gin.H{
			"Playlists": viewModels,
		})
	}
}

func Add(stream *sse.Event, cli *spotifyapi.Client, st store.Store) gin.HandlerFunc {
	return func(c *gin.Context) {
		u := c.PostForm("url")

		urls := strings.Split(u, " ")
		var models []EditListViewModel
		for _, u := range urls {
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

			var tracks []string
			for _, track := range album.Tracks.Tracks {
				tracks = append(tracks, string(track.URI))
			}

			pl := &playlist.Playlist{
				Url:       u,
				Name:      album.Name,
				Img:       imgUrl,
				SpotifyID: spotifyId,
				Artists:   artistNames,
				Tracks:    tracks,
			}

			err = st.Save(c, pl)
			if err != nil {
				_ = c.Error(err)
				return
			}

			model := EditListViewModel{
				ID:        pl.ID.Hex(),
				Img:       pl.Img,
				Name:      pl.Name,
				Artists:   pl.Artists,
				PlayCount: pl.PlayCount,
			}
			models = append(models, model)
		}

		c.HTML(http.StatusOK, "edit-list-entries", gin.H{
			"Playlists": models,
		})

		stream.Message <- sse.UpdateListMessage()
	}
}

func Delete(stream *sse.Event, s store.Store) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		objectId, _ := primitive.ObjectIDFromHex(id)

		err := s.Delete(c, bson.M{"_id": objectId}, nil)
		if err != nil {
			_ = c.Error(err)
			return
		}

		c.Status(http.StatusOK)
		stream.Message <- sse.UpdateListMessage()
	}
}
