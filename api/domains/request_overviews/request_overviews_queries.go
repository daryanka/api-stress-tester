package request_overviews

const (
	queryGetAllRequests = `
SELECT request_overviews.id AS id,
       request_overviews.user_id as user_id,
	   req_name,
       domain_id,
       endpoint,
       method,
       payload,
       time,
       num_requests,
       successful_req,
       failed_req,
       average_response_time,
       status,
	   created_at,
       d.id as nested_domain_id,
       domain_url
FROM request_overviews
    LEFT JOIN domains d on d.id = request_overviews.domain_id
WHERE request_overviews.user_id = ?;`
	queryGetSingle = `
SELECT request_overviews.id AS id,
       request_overviews.user_id as user_id,
	   req_name,
       domain_id,
       endpoint,
       method,
       payload,
       time,
       num_requests,
       successful_req,
       failed_req,
       average_response_time,
       status,
	   created_at,
	   d.id as nested_domain_id,
       domain_url
FROM request_overviews 
	LEFT JOIN domains d on d.id = request_overviews.domain_id
WHERE request_overviews.user_id = ? AND request_overviews.id = ?;`
	queryDelete = `DELETE FROM request_overviews WHERE user_id = ? AND id = ?;`
	queryCreate = `
INSERT INTO request_overviews (
       user_id,
	   req_name,
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
       ?,
       ?);`
	queryUpdateResults = `UPDATE request_overviews SET successful_req = ?, failed_req = ?, average_response_time = ?, status = ? WHERE id = ?;`
	queryUpdateStatus  = `UPDATE request_overviews SET status = ? WHERE id = ?;`
)
