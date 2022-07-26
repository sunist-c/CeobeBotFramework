package http

import (
	"github.com/gin-gonic/gin"
	"github.com/sunist-c/CeobeBotFramework/infrastructure/logging"
)

func init() {
	gin.SetMode(gin.ReleaseMode)
	logging.Debug("init domain/http", "")
}
