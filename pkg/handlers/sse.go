package handlers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

type ClientChan chan string

var clients = make(map[ClientChan]struct{})

func SSE(c *gin.Context) {
	c.Writer.Header().Set("Content-Type", "text/event-stream")
	c.Writer.Header().Set("Cache-Control", "no-cache")
	c.Writer.Header().Set("Connection", "keep-alive")

	log.Log().Msg("client connected")

	eventChan := make(ClientChan)
	clients[eventChan] = struct{}{}

	defer func() {
		delete(clients, eventChan)
		close(eventChan)
	}()

	notify := c.Writer.CloseNotify()
	go func() {
		<-notify
		log.Log().Msg("client disconnected")
	}()

	for {
		data := <-eventChan
		_, _ = fmt.Fprintf(c.Writer, "data: %s\n\n", data)
		c.Writer.Flush()
	}
}

func sendUpdateEvent(data string) {
	for client := range clients {
		client <- data
	}
}
