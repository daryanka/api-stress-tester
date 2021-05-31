package websocket_conn

import (
	"encoding/json"
	"fmt"
	"github.com/daryanka/api-stress-tester/api/domains/request_overviews"
	"github.com/daryanka/api-stress-tester/api/utils"
	"github.com/gorilla/websocket"
	"sync"
)

type webSocketData struct {
	ID   int64
	Data []byte
}

type HubT struct {
	mu            sync.Mutex
	Connections   map[int64]*websocket.Conn
	WriteChannels map[int64]chan webSocketData

	RequestDoneChan map[int64]chan struct{}
}

type UserMessage struct {
	Type      string `json:"type"`
	RequestID int64  `json:"request_id"`
}

var Hub HubT

func NewHub() {
	Hub = HubT{
		Connections:   make(map[int64]*websocket.Conn),
		WriteChannels: map[int64]chan webSocketData{},

		RequestDoneChan: make(map[int64]chan struct{}),
	}
}

func HandleAddConnection(conn *websocket.Conn, userID int64) {
	// Check if user already exists
	Hub.mu.Lock()

	if _, exists := Hub.Connections[userID]; exists {
		// Close existing connection
		closeConnection(userID)
	}

	// Add connection
	Hub.Connections[userID] = conn

	// Create new channel
	ch := make(chan webSocketData, 20)
	Hub.WriteChannels[userID] = ch
	Hub.mu.Unlock()

	go handleConn(conn, userID)
	go handleSendMessage(ch, userID)
}

// READ MESSAGES
func handleConn(conn *websocket.Conn, id int64) {
	for {
		mt, msg, err := conn.ReadMessage()
		if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
			closeConnection(id)
			return
		}
		if mt == -1 {
			closeConnection(id)
			break
		}
		// Handle Messages
		fmt.Println("message: ", string(msg))
		message := UserMessage{}
		if e := json.Unmarshal(msg, &message); e != nil {
			utils.Logger.Error("error invalid user message JSON ", err)
		}

		if message.Type == "CANCEL_REQUEST" {
			CancelRequest(message.RequestID, id)
		}
	}
}

func WriteMessage(id int64, data []byte) {
	Hub.mu.Lock()
	if ch, ok := Hub.WriteChannels[id]; ok {
		ch <- webSocketData{
			ID:   id,
			Data: data,
		}
	}
	Hub.mu.Unlock()
}

// WRITE MESSAGES
func handleSendMessage(ch <-chan webSocketData, id int64) {
	for {
		msg := <-ch
		if string(msg.Data) == "CLOSE" {
			break
		}
		Hub.mu.Lock()
		if val, ok := Hub.Connections[msg.ID]; ok {
			// messageType 1 = TextMessage
			Hub.mu.Unlock()
			err := val.WriteMessage(1, msg.Data)
			if err != nil {
				closeConnection(id)
				utils.Logger.Error("error writing message, connection closed ", err)
			}
		} else {
			Hub.mu.Unlock()
		}
	}
}

func closeConnection(id int64) {
	// Close + Remove Websocket Connection
	Hub.mu.Lock()
	defer Hub.mu.Unlock()
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

func AddDoneChan(reqID int64, ch chan struct{}) {
	Hub.mu.Lock()
	Hub.RequestDoneChan[reqID] = ch
	Hub.mu.Unlock()
}

func RemoveDoneChan(reqID int64) {
	Hub.mu.Lock()

	if _, ok := Hub.RequestDoneChan[reqID]; ok {
		// Remove
		delete(Hub.RequestDoneChan, reqID)
	}

	Hub.mu.Unlock()
}

func CancelRequest(requestID, userID int64) {
	// Check if user is allowed to cancel request
	_, err := request_overviews.RequestOverviewDao.GetSingle(userID, requestID)
	if err != nil {
		return
	}

	Hub.mu.Lock()
	if val, ok := Hub.RequestDoneChan[requestID]; ok {
		close(val)
		delete(Hub.RequestDoneChan, requestID)
	}
	Hub.mu.Unlock()
}
