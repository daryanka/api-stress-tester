package users

import (
	"database/sql"
	"github.com/daryanka/api-stress-tester/api/clients"
	"github.com/daryanka/api-stress-tester/api/utils"
	"os"
)

type UserDoaI interface {
	Find(id int64) (u *User, err error)
	Create(u *User) (id int64, err error)
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
