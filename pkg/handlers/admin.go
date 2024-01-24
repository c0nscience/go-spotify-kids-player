package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	spotifyapi "github.com/zmb3/spotify/v2"
	"go-spotify-kids-player/pkg/playlist"
	"go-spotify-kids-player/pkg/spotify"
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

func Add(cli *spotifyapi.Client, st store.Store) gin.HandlerFunc {
	return func(c *gin.Context) {
		u := c.PostForm("url")

		log.Info().Msgf("value: %v", u)

		urls := strings.Split(u, " ")
		var models []EditListViewModel
		for _, u := range urls {
			log.Info().Msgf("line: %s", u)
			spotifyId := spotify.GetIdFrom(u)
			log.Info().Msgf("spotifyId: %s", spotifyId)

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

			pl := &playlist.Playlist{
				Url:       u,
				Name:      album.Name,
				Img:       imgUrl,
				SpotifyID: spotifyId,
				Artists:   artistNames,
			}

			id, err := st.Save(c, pl)
			if err != nil {
				_ = c.Error(err)
				return
			}
			pl.ID = id.(primitive.ObjectID)

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

		sendUpdateEvent("")
	}
}

func Delete(s store.Store) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		objectId, _ := primitive.ObjectIDFromHex(id)

		err := s.Delete(c, bson.M{"_id": objectId}, nil)
		if err != nil {
			_ = c.Error(err)
			return
		}

		c.Status(http.StatusOK)
		sendUpdateEvent("")
	}
}
