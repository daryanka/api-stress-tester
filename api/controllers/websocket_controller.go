package controllers

import (
	"fmt"
	"github.com/daryanka/api-stress-tester/api/utils"
	"github.com/daryanka/api-stress-tester/api/websocket_conn"
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

	u := GetAuthUser(c)

	conn, err := w.Upgrade(c.Writer, c.Request, http.Header{
		"Sec-WebSocket-Protocol": []string{c.GetHeader("Sec-WebSocket-Protocol")},
	})
	if err != nil {
		fmt.Println(err.Error())
		e := utils.StandardInternalServerError()
		c.JSON(e.Code(), e)
		return
	}

	go websocket_conn.HandleAddConnection(conn, u.ID)
}

func handleConn(conn *websocket.Conn) {
	for {
		mt, message, err := conn.ReadMessage()
		fmt.Println("messageType: ", mt)
		fmt.Println("message: ", string(message))

		if err != nil {
			if mt == -1 {
				// TODO  close websocket / remove from existing
				break
			}
			utils.Logger.Error("error reading message", err)
			break
		}
		err = conn.WriteMessage(mt, []byte("bye world"))
		if err != nil {
			utils.Logger.Error("error writing message", err)
			break
		}
	}
}
