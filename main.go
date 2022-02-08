package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/joooostb/wahoo-garmin-sync/pkg/dropbox"
	log "github.com/sirupsen/logrus"
)

func main() {
	g := gin.Default()
	// Only log the warning severity or above.
	log.SetLevel(log.DebugLevel)

	g.GET("/healthz", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})
	g.GET("/", func(c *gin.Context) {
		dropbox.Authorize(c)
	})
	g.GET("/oauth2", func(c *gin.Context) {
		dropbox.Authenticate(c)
	})
	g.GET("/webhook", func(c *gin.Context) {
		dropbox.Challenge(c)
	})
	g.POST("/webhook", func(c *gin.Context) {
		dropbox.Handler(c)
	})
	g.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
