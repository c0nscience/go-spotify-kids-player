package spotify

import (
	"context"
	"github.com/zmb3/spotify/v2"
)

func GetCoverImage(ctx context.Context, cli *spotify.Client, id string) (string, error) {
	album, err := cli.GetAlbum(ctx, spotify.ID(id))
	if err != nil {
		return "", err
	}
	return album.Images[0].URL, nil
}
