package playlist

import (
	"go-spotify-kids-player/pkg/store"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var Collection store.CollectionName = "playlist"

type Playlist struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	Name      string             `bson:"name"`
	Url       string             `bson:"url"`
	Img       string             `bson:"img"`
	SpotifyID string             `bson:"spotifyID"`
	Artists   []string           `bson:"artists"`
}

func (p *Playlist) ObjectId() primitive.ObjectID {
	return p.ID
}

func (p *Playlist) SetObjectId(id primitive.ObjectID) {
	p.ID = id
}

var _ store.HasObjectId = &Playlist{}
