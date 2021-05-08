package user

import (
	"database/sql"
	"github.com/daryanka/api-stress-tester/api/clients"
	"github.com/daryanka/api-stress-tester/api/utils"
	"os"
)

type UserDoaI interface {
	Find(id int64) (u *User, err error)
	Create(u *User) (id int64, err error)
	FindByEmail(email string) (u User, err error)
	EmailInUse(email string) (bool, error)
	VerifyEmail(token string) (bool, error)
}

type userDao struct{}

var UserDao UserDoaI = &userDao{}

func (i *userDao) Find(id int64) (u *User, err error) {
	key := os.Getenv("ENC_KEY")

	err = clients.DB.Get(&u, queryGetUser, key, key, id)
	if err != nil && err != sql.ErrNoRows {
		utils.Logger.Error("error fetching user", err)
	}
	return
}

func (i *userDao) FindByEmail(email string) (u User, err error) {
	key := os.Getenv("ENC_KEY")

	err = clients.DB.Get(&u, queryGetUserByEmail, key, key, key, email)
	if err != nil && err != sql.ErrNoRows {
		utils.Logger.Error("error fetching user by email", err)
	}
	return
}

func (i *userDao) Create(u *User) (id int64, err error) {
	key := os.Getenv("ENC_KEY")

	res, err := clients.DB.Exec(queryCreateUser,
		u.Name, key,
		u.Email, key,
		u.EmailToken,
		u.Password,
	)

	if err != nil {
		utils.Logger.Error("error creating user", err)
		return
	}

	id, err = res.LastInsertId()
	if err != nil {
		utils.Logger.Error("error getting created user id", err)
	}
	return
}

func (i *userDao) EmailInUse(email string) (bool, error) {
	var id int64
	key := os.Getenv("ENC_KEY")

	err := clients.DB.Get(&id, queryEmailInUse, key, email)
	if err != nil && err != sql.ErrNoRows {
		utils.Logger.Error("error getting email in use", err)
		return false, err
	}

	if err == sql.ErrNoRows {
		return false, nil
	}

	return true, nil
}

func (i *userDao) VerifyEmail(token string) (bool, error) {
	res, err := clients.DB.Exec(queryUpdateEmailVerified, 1, token)
	if err != nil {
		utils.Logger.Error("error updating email_verified", err)
		return false, err
	}

	num, err := res.RowsAffected()
	if err != nil {
		utils.Logger.Error("error getting rows affected email_verified", err)
		return false, err
	}

	if num == 0 {
		return false, nil
	}

	return true, nil
}