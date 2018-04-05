package api

import (
	"github.com/ONSBR/Plataforma-EventManager/bus"
	"github.com/ONSBR/Plataforma-EventManager/infra"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

//BuildAPI register all endpoints that should be exposed
func BuildAPI(broker bus.Dispatcher) {
	port := infra.GetEnv("PORT", "8081")
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	log.Info("Registering query routes")
	registerQueryEndpoints(r)
	log.Info("Registering commands routes")
	registerCommandsApi(r, broker)
	log.Info("Listening on: 0.0.0.0:" + port)
	r.Run(":" + port)
}
