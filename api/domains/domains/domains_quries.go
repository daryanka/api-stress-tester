package domains

const (
	queryGetAll = `
SELECT d.id AS id,
       domain_url,
       token,
       verified,
       d.user_id AS user_id,
       COALESCE(COUNT(ro.id), 0) AS num_requests
FROM domains d
         LEFT JOIN request_overviews ro on d.id = ro.domain_id
WHERE d.user_id = ? GROUP BY d.id;
`
	queryGetSingle = `
SELECT d.id AS id,
       domain_url,
       token,
       verified,
       d.user_id AS user_id,
       COALESCE(COUNT(ro.id), 0) AS num_requests
FROM domains d
         LEFT JOIN request_overviews ro on d.id = ro.domain_id
WHERE d.user_id = ? AND d.id = ? GROUP BY d.id;
`
	queryAddDomain = `INSERT INTO domains (domain_url, token, verified, user_id) VALUES (?, ?, ?, ?);`
	queryDelete = `DELETE FROM domains WHERE id = ? AND user_id = ?;`
	queryVerifyDomain = `UPDATE domains SET verified = 1 WHERE id = ?;`
)