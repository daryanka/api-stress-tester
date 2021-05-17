package domains

const (
	queryGetAll = `
SELECT d.id AS id,
       domain_url,
       token,
       verified,
       user_id,
       COALESCE(COUNT(ro.id), 0) AS num_requests
FROM domains d
         LEFT JOIN request_overviews ro on d.id = ro.domain_id
WHERE status != 0
  AND user_id = ?;
`
	queryAddDomain = `INSERT INTO domains (domain_url, token, verified, user_id) VALUES (?, ?, ?, ?);`
	queryDelete = `DELETE FROM domains WHERE id = ? AND user_id = ?;`
	queryVerifyDomain = `UPDATE domains SET verified = 1 WHERE id = ?;`
)