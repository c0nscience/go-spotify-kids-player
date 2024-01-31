package ha

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"github.com/rs/zerolog/log"
	"sync/atomic"
)

var msgId atomic.Int32

func Handle(msgChan chan interface{}) {
	msgId.Store(0)
	go func() {
		for {
			select {
			case m := <-msgChan:
				msg := m.(map[string]interface{})
				log.Info().Msgf("handle: %v", msg)
				switch msg["type"].(string) {
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
				}
			}
		}
	}()
}
