package api

import (
	"bytes"
	"io/ioutil"
	"time"

	"github.com/ONSBR/Plataforma-EventManager/domain"
	"github.com/ONSBR/Plataforma-EventManager/flow"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func registerCommandsAPI(r *gin.Engine) {
	fullEventFlow := flow.GetEventRouter()
	storeEventFlow := flow.GetBasicEventRouter()
	r.Use(func(c *gin.Context) {
		var bodyBytes []byte
		if c.Request.Body != nil {
			bodyBytes, _ = ioutil.ReadAll(c.Request.Body)
		}
		// Restore the io.ReadCloser to its original state
		c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
		log.Info(string(bodyBytes))
		c.Next()
	})
	r.PUT("/sendevent", func(c *gin.Context) {
		event := domain.NewEvent()
		var err error
		if err = c.BindJSON(event); err == nil {

			begin := time.Now()
			err = fullEventFlow.Push(event)
			if err != nil {
				log.Error(err)
			}
			log.Info("Tempo total de processamento do evento:", time.Now().Sub(begin))
			if err == nil {
				c.JSON(200, gin.H{
					"message": "OK",
				})
				return
			}
		}
		log.Error(err.Error())
		c.JSON(400, gin.H{
			"message": err.Error(),
		})
	})
	log.Debug("Register route POST /save")
	r.POST("/save", func(c *gin.Context) {
		event := new(domain.Event)
		if err := c.BindJSON(event); err != nil {
			c.JSON(400, gin.H{
				"message": err.Error(),
			})
		}

		begin := time.Now()
		if err := storeEventFlow.Push(event); err != nil {
			log.Error(err.Error())
			c.JSON(200, err)
		} else {
			c.JSON(200, event)
		}
		log.Info("Tempo total de processamento do evento:", time.Now().Sub(begin))
	})
}
