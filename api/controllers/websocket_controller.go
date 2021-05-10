package controllers

import (
	"fmt"
	"github.com/daryanka/api-stress-tester/api/utils"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"net/http"
)

type WebSocketControllerI interface {
	Connect(c *gin.Context)
}

type webSocketController struct{}

var WebSocketController WebSocketControllerI = &webSocketController{}

func (i *webSocketController) Connect(c *gin.Context) {
	w := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	h := c.GetHeader("Sec-WebSocket-Protocol")

	conn, err := w.Upgrade(c.Writer, c.Request, http.Header{
		"Sec-WebSocket-Protocol": []string{
			h,
		},
	})
	if err != nil {
		fmt.Println(err.Error())
		e := utils.StandardInternalServerError()
		c.JSON(e.Code(), e)
		return
	}
	r := conn.Subprotocol()
	fmt.Print("sub proto call is ", r)
}