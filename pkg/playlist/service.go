package playlist

import "go-spotify-kids-player/pkg/models"

var store = map[string]models.Playlist{
	"1": {Id: "1", Url: "https://open.spotify.com/album/60r4JW9AcLzEPrQIRoJ31p"},
	"2": {Id: "2", Url: "https://open.spotify.com/album/7lsqXCAIcZNVVYB1e8dRad?si=FMx7hyhWRtiOvRWRPZTW4w"},
}

func GetAll() []models.Playlist {
	res := []models.Playlist{}

	for _, playlist := range store {
		res = append(res, playlist)
	}

	return res
}

func GetById(id string) models.Playlist {
	return store[id]
}
