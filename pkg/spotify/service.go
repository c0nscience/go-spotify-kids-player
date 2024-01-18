package spotify

import (
	"context"
	spotifyapi "github.com/zmb3/spotify/v2"
)

func GetAlbum(ctx context.Context, cli *spotifyapi.Client, id string) (*spotifyapi.FullAlbum, error) {
	album, err := cli.GetAlbum(ctx, spotifyapi.ID(id))
	if err != nil {
		return nil, err
	}
	return album, nil
}
