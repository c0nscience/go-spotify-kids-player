package ha

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/gorilla/websocket"
	"github.com/rs/zerolog/log"
	"go-spotify-kids-player/pkg/player"
	"go-spotify-kids-player/pkg/playlist"
	"go-spotify-kids-player/pkg/sse"
	"go-spotify-kids-player/pkg/store"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"strings"
	"sync/atomic"
)

var msgId atomic.Int32

func Handle(stream *sse.Event, store store.Store, msgChan chan interface{}) chan struct{} {
	msgId.Store(0)
	connectedChan := make(chan struct{})
	go func() {
		for {
			select {
			case m := <-msgChan:
				msg := m.(map[string]interface{})
				log.Info().Msgf("handle: %v", msg)
				switch msg["type"].(string) {
				case "result":
					if msg["success"].(bool) {
						close(connectedChan)
					}
				case "auth_required":
					b, err := json.Marshal(&map[string]string{
						"type":         "auth",
						"access_token": accessToken,
					})
					if err != nil {
						log.Error().Err(err)
					}
					err = WebsocketConnection.WriteMessage(websocket.TextMessage, b)
					if err != nil {
						log.Error().Err(err)
					}
				case "auth_ok":
					b, err := json.Marshal(&map[string]interface{}{
						"id":         msgId.Add(1),
						"type":       "subscribe_events",
						"event_type": "state_changed",
					})
					if err != nil {
						log.Error().Err(err)
					}
					err = WebsocketConnection.WriteMessage(websocket.TextMessage, b)
					if err != nil {
						log.Error().Err(err)
					}
				case "event":
					e, _ := event(msg)
					d, _ := data(e)

					id, _ := entityId(d)
					if !strings.HasPrefix(id, "media_player") {
						continue
					}

					ns, _ := newState(d)
					attr, _ := attributes(ns)
					contentId, _ := mediaContentId(attr)

					var pl playlist.Playlist
					err := store.FindOne(context.Background(), bson.M{"tracks": contentId}, &pl)
					if err != nil {
						if !errors.Is(err, mongo.ErrNoDocuments) {
							log.Error().Err(err).Msg("could not find a playlist with this track")
						}
						continue
					}

					s, _ := state(ns)
					log.Info().Msgf("found: %s %s %s on %s", pl.Name, pl.ID.Hex(), s, id)

					currentState := player.State()
					if s == "playing" && !currentState.Playing {
						player.Reduce(player.Action{
							Type:    player.PlayType,
							Payload: []string{pl.ID.Hex(), pl.Img},
						})
						stream.Message <- sse.UpdatePlayerMessage()
					} else if s == "paused" && currentState.Playing {
						player.Reduce(player.Action{
							Type: player.StopType,
						})
						stream.Message <- sse.UpdatePlayerMessage()
					}

				}
			}
		}
	}()

	return connectedChan
}
