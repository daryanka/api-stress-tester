package user

const (
	queryCreateUser = `INSERT INTO users (name, email, email_token, password) VALUES (AES_ENCRYPT(?, ?), AES_ENCRYPT(?, ?), ?, ?);`
	queryGetUser    = `SELECT 
	id,	
	AES_DECRYPT(name, ?) AS name,
	AES_DECRYPT(email, ?) AS email,
	email_token,
	password,
	email_verified
FROM users WHERE id = ?;`
	queryGetUserByEmail    = `SELECT 
	id,	
	AES_DECRYPT(name, ?) AS name,
	AES_DECRYPT(email, ?) AS email,
	email_token,
	password,
	email_verified
FROM users WHERE AES_DECRYPT(email, ?) = ?;`
	queryEmailInUse = `SELECT id FROM users WHERE AES_DECRYPT(email, ?) = ?;`
	queryUpdateEmailVerified = `UPDATE users SET email_verified = ? WHERE email_token = ?;`
)
