package middleware

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

func ZerologMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		latencyMs := float64(time.Since(start).Microseconds()) / 1000.0

		log.Info().
			Str("method", c.Request.Method).
			Str("path", c.Request.URL.Path).
			Any("query", c.Request.URL.Query()).
			Int("status", c.Writer.Status()).
			Str("latency", fmt.Sprintf("%.3f ms", latencyMs)).
			Str("client_ip", c.ClientIP()).
			Msg("request")
	}
}
