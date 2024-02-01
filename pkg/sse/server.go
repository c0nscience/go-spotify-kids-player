package sse

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

type Event struct {
	Message       chan Message
	NewClients    chan ClientChan
	ClosedClients chan ClientChan
	TotalClients  map[ClientChan]bool
}

type Message struct {
	Name    string
	Content string
}

type ClientChan chan Message

func NewServer() (event *Event) {
	event = &Event{
		Message:       make(chan Message),
		NewClients:    make(chan ClientChan),
		ClosedClients: make(chan ClientChan),
		TotalClients:  make(map[ClientChan]bool),
	}

	go event.Listen()

	return
}

func (stream *Event) Listen() {
	for {
		select {
		case client := <-stream.NewClients:
			stream.TotalClients[client] = true
			log.Info().Msgf("Client added. %d registered clients", len(stream.TotalClients))

		case client := <-stream.ClosedClients:
			delete(stream.TotalClients, client)
			close(client)
			log.Info().Msgf("Removed client. %d registered clients", len(stream.TotalClients))

		case eventMsg := <-stream.Message:
			for clientMessageChan := range stream.TotalClients {
				clientMessageChan <- eventMsg
			}
		}
	}
}

func (stream *Event) ServeHTTP() gin.HandlerFunc {
	return func(c *gin.Context) {
		clientChan := make(ClientChan)

		stream.NewClients <- clientChan

		defer func() {
			stream.ClosedClients <- clientChan
		}()

		c.Set("clientChan", clientChan)

		c.Next()
	}
}

func HeadersMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Content-Type", "text/event-stream")
		c.Writer.Header().Set("Cache-Control", "no-cache")
		c.Writer.Header().Set("Connection", "keep-alive")
		c.Writer.Header().Set("Transfer-Encoding", "chunked")
		c.Next()
	}
}
func UpdateListMessage() Message {
	return Message{
		Name:    "update-list",
		Content: "",
	}
}
func UpdatePlayerMessage() Message {
	return Message{
		Name:    "update-player",
		Content: "",
	}
}
