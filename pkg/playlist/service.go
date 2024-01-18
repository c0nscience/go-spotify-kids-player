package playlist

import (
	"github.com/google/uuid"
)

type Playlist struct {
	ID        string
	Name      string
	Url       string
	Img       string
	SpotifyID string
	Artists   []string
}

var store = map[string]Playlist{}

func GetAll() []Playlist {
	var res []Playlist

	for _, playlist := range store {
		res = append(res, playlist)
	}

	return res
}

func GetById(id string) Playlist {
	return store[id]
}

func Store(p Playlist) Playlist {
	id, _ := uuid.NewRandom()
	p.ID = id.String()

	store[id.String()] = p

	return p
}

func DeleteByID(id string) {
	delete(store, id)
}
