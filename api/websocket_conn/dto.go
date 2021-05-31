package websocket_conn

const (
	INDIVIDUAL_REQUEST = "INDIVIDUAL_REQUEST"
	REQUEST_FAILED = "REQUEST_FAILED"
	REQUEST_COMPLETE = "REQUEST_COMPLETE"
)

type MessageType struct {
	Type   string      `json:"type"`
	Status int         `json:"status"`
	Info   interface{} `json:"info"`
}

type RequestUpdate struct {
	RequestID int64  `json:"request_id"`
	Message   string `json:"message"`
}

type CompletedIndividualRequest struct {
	RequestID  int64 `json:"request_id"`
	TimeTaken  int64 `json:"time_taken"`
	StatusCode int   `json:"status_code"`
}
