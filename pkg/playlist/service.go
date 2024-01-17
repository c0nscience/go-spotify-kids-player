package playlist

import (
	"net/url"
	"strings"
)

type Playlist struct {
	Id   string
	Url  string
	Name string
	Img  string
}

func GetId(playlist Playlist) string {
	u, _ := url.Parse(playlist.Url)
	p := strings.Split(u.Path[1:], "/")
	return p[len(p)-1]
}

var store = map[string]Playlist{
	"1": {
		Id:   "1",
		Name: "Raeuber",
		Url:  "https://open.spotify.com/album/60r4JW9AcLzEPrQIRoJ31p",
		Img:  "https://i.scdn.co/image/ab67616d0000b2739c972d46a118bc9a0add63e7",
	},
	"2": {
		Id:   "2",
		Name: "Hexe",
		Url:  "https://open.spotify.com/album/7lsqXCAIcZNVVYB1e8dRad?si=FMx7hyhWRtiOvRWRPZTW4w",
		Img:  "https://i.scdn.co/image/ab67616d0000b2739c972d46a118bc9a0add63e7",
	},
}

func GetAll() []Playlist {
	res := []Playlist{}

	for _, playlist := range store {
		res = append(res, playlist)
	}

	return res
}

func GetById(id string) Playlist {
	return store[id]
}
