package api

import (
	"github.com/ONSBR/Plataforma-EventManager/bus"
	"github.com/ONSBR/Plataforma-EventManager/domain"
	"github.com/ONSBR/Plataforma-EventManager/flow"
	"github.com/ONSBR/Plataforma-EventManager/infra"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func registerCommandsApi(r *gin.Engine, dispatcher bus.Dispatcher) {
	fullEventFlow := flow.GetEventFlow(dispatcher)
	storeEventFlow := flow.GetStoreEventFlow(dispatcher)
	log.Info("Register route PUT /sendevent")
	r.PUT("/sendevent", func(c *gin.Context) {
		event := domain.NewEvent()
		if err := c.BindJSON(event); err != nil {
			c.JSON(400, gin.H{
				"message": err.Error(),
			})
		} else if err := fullEventFlow.Push(event); err != nil {
			ex := err.(*infra.Exception)
			c.JSON(ex.HTTPStatus(), ex)
		} else {
			c.JSON(200, gin.H{
				"message": "OK",
			})
		}
	})
	log.Info("Register route POST /save")
	r.POST("/save", func(c *gin.Context) {
		event := new(domain.Event)
		if err := c.BindJSON(event); err != nil {
			c.JSON(400, gin.H{
				"message": err.Error(),
			})
		} else if err := storeEventFlow.Push(event); err != nil {
			ex := err.(*infra.Exception)
			c.JSON(ex.HTTPStatus(), ex)
		} else {
			c.JSON(200, gin.H{
				"message": "OK",
			})
		}
	})
}
