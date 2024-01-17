package models

type PlaylistViewModel struct {
	Img string
	Id  string
}

type ListViewModel struct {
	Playlists []PlaylistViewModel
}
