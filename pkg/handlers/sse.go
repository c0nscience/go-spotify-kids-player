package handlers

import (
	"github.com/gin-gonic/gin"
	"go-spotify-kids-player/pkg/sse"
	"io"
)

func SSE(c *gin.Context) {
	v, ok := c.Get("clientChan")
	if !ok {
		return
	}
	clientChan, ok := v.(sse.ClientChan)
	if !ok {
		return
	}
	c.Stream(func(w io.Writer) bool {
		if msg, ok := <-clientChan; ok {
			c.SSEvent(msg.Name, msg.Content)
			return true
		}
		return false
	})
}
