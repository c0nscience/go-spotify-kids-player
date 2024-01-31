package ha

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/rs/zerolog/log"
	"net/url"
)

var WebsocketConnection *websocket.Conn

func Listen() chan interface{} {
	u, err := url.Parse(fmt.Sprintf("%s/api/websocket", host))
	if err != nil {
		log.Error().Err(err).Msg("could not create websocket url")
	}
	u.Scheme = "ws"

	WebsocketConnection, _, err = websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Error().Err(err).Msg("could not connect to websocket")
		return nil
	}

	msgChan := make(chan interface{})

	go func() {
		for {
			_, m, err := WebsocketConnection.ReadMessage()
			if err != nil {
				log.Error().Err(err)
				return
			}

			var msg map[string]interface{}
			err = json.Unmarshal(m, &msg)
			if err != nil {
				log.Error().Err(err)
				return
			}

			msgChan <- msg

			log.Info().Msgf("received: %v", msg)
		}
	}()

	return msgChan
}
