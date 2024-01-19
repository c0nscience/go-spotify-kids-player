package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"net/http"
)

func ErrorHandling(c *gin.Context) {
	c.Next()

	lastErr := c.Errors.Last()
	if lastErr == nil {
		return
	}

	for _, err := range c.Errors {
		log.Error().Err(err)
	}

	c.AbortWithStatus(http.StatusInternalServerError)
}
