package domains

const (
	queryGetAll = `SELECT id, domain_url, token, verified, user_id FROM domains WHERE user_id = ?;`
	queryAddDomain = `INSERT INTO domains (domain_url, token, verified, user_id) VALUES (?, ?, ?, ?);`
	queryDelete = `DELETE FROM domains WHERE id = ? AND user_id = ?;`
	queryVerifyDomain = `UPDATE domains SET verified = 1 WHERE id = ?;`
)