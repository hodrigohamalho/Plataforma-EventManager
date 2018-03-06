package main

import (
	"os"

	"github.com/ONSBR/Plataforma-EventManager/infra"

	"github.com/ONSBR/Plataforma-EventManager/actions"
	"github.com/ONSBR/Plataforma-EventManager/domain"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func init() {
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)
}

func main() {
	log.Info("Starting Event Manager")
	port := infra.GetEnv("PORT", "8081")
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	r.PUT("/sendevent", func(c *gin.Context) {
		log.Info("Pushing event to executor")
		event := new(domain.Event)
		if err := c.BindJSON(event); err != nil {
			c.JSON(400, gin.H{
				"message": err.Error(),
			})
		} else if err := actions.PushEventToExecutor(*event); err != nil {
			c.JSON(400, gin.H{
				"message": err.Error(),
			})
		} else {
			c.JSON(200, gin.H{
				"message": "OK",
			})
		}
	})
	log.Info("Listening on: 0.0.0.0:" + port)
	r.Run(":" + port)
}
