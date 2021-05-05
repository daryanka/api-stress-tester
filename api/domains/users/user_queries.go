package users

const (
	queryCreateUser = `INSERT INTO users (name, email, email_token, password) VALUES (AES_ENCRYPT(?, ?), AES_ENCRYPT(?, ?), ?, ?);'`
	queryGetUser    = `SELECT 
	id,	
	AES_DECRYPT(name, ?) AS name,
	AES_DECRYPT(email, ?) AS email,
	email_token,
	password
FROM users WHERE id = ?;`
)
