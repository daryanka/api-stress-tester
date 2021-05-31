package request_overviews

import (
	"github.com/daryanka/api-stress-tester/api/domains/individual_requests"
	"time"
)

const (
	StatusPending = iota
	StatusInProgress
	StatusComplete
	StatusCancelled

)
type RequestOverview struct {
	ID                  int64                                   `json:"id" db:"id"`
	UserID              int64                                   `json:"user_id" db:"user_id"`
	DomainID            int64                                   `json:"domain_id" db:"domain_id"`
	Endpoint            string                                  `json:"endpoint" db:"endpoint"`
	Method              string                                  `json:"method" db:"method"`
	Payload             *string                                 `json:"payload" db:"payload"`
	Time                int                                     `json:"time" db:"time"`
	NumRequests         int                                     `json:"num_requests" db:"num_requests"`
	SuccessfulReq       int                                     `json:"successful_req" db:"successful_req"`
	FailedReq           int                                     `json:"failed_req" db:"failed_req"`
	AverageResponseTime float64                                 `json:"average_response_time" db:"average_response_time"`
	Status              int                                     `json:"status" db:"status"`
	CreatedAt           time.Time                               `json:"created_at" db:"created_at"`
	IndividualRequests  []individual_requests.IndividualRequest `json:"individual_requests"`
}

type NewRequest struct {
	UserID      int64
	Payload     *string `json:"payload"`
	DomainID    int64   `json:"domain_id" validate:"required"`
	Endpoint    string  `json:"endpoint" validate:"required"`
	Method      string  `json:"method" validate:"required,oneof=GET HEAD POST PUT PATCH DELETE CONNECT OPTIONS TRACE"`
	Time        int     `json:"time" validate:"required,max=5,min=0"`
	NumRequests int     `json:"num_requests" validate:"required"`
}
