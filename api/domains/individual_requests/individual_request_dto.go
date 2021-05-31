package individual_requests

import "time"

type IndividualRequest struct {
	ID                int64     `json:"id" db:"id"`
	RequestOverviewId int64     `json:"request_overview_id" db:"request_overview_id"`
	StatusCode        int       `json:"status_code" db:"status_code"`
	TimeTaken         int64     `json:"time_taken" db:"time_taken"`
	CreatedAt         time.Time `json:"created_at" db:"created_at"`
}
