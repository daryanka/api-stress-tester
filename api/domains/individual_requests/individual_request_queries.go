package individual_requests

const (
	queryGetAll = `
SELECT id,
       request_overview_id,
       status_code,
       time_taken,
	   created_at
FROM individual_requests
WHERE request_overview_id = ?;`
	queryCreate = `
INSERT INTO individual_requests (request_overview_id, status_code, time_taken) VALUES `
)