package models

import (
	"net/url"
	"strings"
)

type Playlist struct {
	Id  string
	Url string
}

func GetId(playlist Playlist) string {
	u, _ := url.Parse(playlist.Url)
	p := strings.Split(u.Path[1:], "/")
	return p[len(p)-1]
}
