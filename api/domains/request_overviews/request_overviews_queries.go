package request_overviews

const (
	queryGetAllRequests = `
SELECT id,
       user_id,
       domain_id,
       endpoint,
       method,
       payload,
       time,
       num_requests,
       successful_req,
       failed_req,
       average_response_time,
       status
FROM request_overviews WHERE user_id = ?;`
	queryGetSingle = `
SELECT id,
       user_id,
       domain_id,
       endpoint,
       method,
       payload,
       time,
       num_requests,
       successful_req,
       failed_req,
       average_response_time,
       status
FROM request_overviews WHERE user_id = ? AND id = ?;`
	queryDelete = `DELETE FROM request_overviews WHERE user_id = ? AND id = ?;`
	queryCreate = `
INSERT INTO request_overviews (
       user_id,
       domain_id,
       endpoint,
       method,
       payload,
       time,
       num_requests,
       successful_req,
       failed_req,
       average_response_time,
       status) VALUES (
	   ?,
       ?,
       ?,
       ?,
       ?,
       ?,
       ?,
       ?,
       ?,
       ?,
       ?);`
	queryUpdateResults= `UPDATE request_overviews SET successful_req = ?, failed_req = ?, average_response_time = ?, status = ? WHERE id = ?;`
)