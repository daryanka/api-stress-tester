package websocket_conn

import (
	"fmt"
	"github.com/daryanka/api-stress-tester/api/utils"
	"github.com/gorilla/websocket"
)

type webSocketData struct {
	ID   int64
	Data []byte
}

type HubT struct {
	Connections map[int64]*websocket.Conn

	WriteChannels map[int64]chan webSocketData
}

var Hub HubT

func NewHub() {
	Hub = HubT{
		Connections:   make(map[int64]*websocket.Conn),
		WriteChannels: map[int64]chan webSocketData{},
	}
}

func HandleAddConnection(conn *websocket.Conn, id int64) {
	// Check if user already exists
	if _, exists := Hub.Connections[id]; exists {
		// Close existing connection
		closeConnection(id)
	}

	// Add connection
	Hub.Connections[id] = conn

	// Create new channel
	ch := make(chan webSocketData, 20)
	Hub.WriteChannels[id] = ch

	go handleConn(conn, id)
	go handleSendMessage(ch, id)
}

// READ MESSAGES
func handleConn(conn *websocket.Conn, id int64) {
	for {
		fmt.Println("test")
		mt, _, err := conn.ReadMessage()
		if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
			closeConnection(id)
			return
		}
		if mt == -1 {
			closeConnection(id)
			break
		}
	}
}

func WriteMessage(id int64, data []byte) {
	if ch, ok := Hub.WriteChannels[id]; ok {
		ch <- webSocketData{
			ID:   id,
			Data: data,
		}
	}
}

// WRITE MESSAGES
func handleSendMessage(ch <-chan webSocketData, id int64) {
	for {
		msg := <-ch
		if string(msg.Data) == "CLOSE" {
			break
		}
		if val, ok := Hub.Connections[msg.ID]; ok {
			// messageType 1 = TextMessage
			err := val.WriteMessage(1, msg.Data)
			if err != nil {
				closeConnection(id)
				utils.Logger.Error("error writing message, connection closed ", err)
			}
		}
	}
}

func closeConnection(id int64) {
	// Close + Remove Websocket Connection
	if val, ok := Hub.Connections[id]; ok {
		err := val.Close()
		if err != nil {
			utils.Logger.Error("error closing connection", err)
		}

		// Remove connection
		delete(Hub.Connections, id)
	}

	// Close + Remove Channel
	if val, ok := Hub.WriteChannels[id]; ok {
		val <- webSocketData{
			ID:   id,
			Data: []byte(`CLOSE`),
		}
		close(val)
		delete(Hub.WriteChannels, id)
	}
}
