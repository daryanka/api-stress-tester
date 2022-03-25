package services

import (
	"database/sql"
	"github.com/daryanka/api-stress-tester/api/domains/user"
	"github.com/daryanka/api-stress-tester/api/utils"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserServiceI interface {
	Login(req *user.ReqLogin) (*utils.TokenWithClaims, utils.RestErrI)
	Register(req *user.ReqRegister) utils.RestErrI
}

type userService struct{}

var UserService UserServiceI = &userService{}

func (i *userService) Login(req *user.ReqLogin) (*utils.TokenWithClaims, utils.RestErrI) {
	foundUser, err := user.UserDao.FindByEmail(req.Email)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, utils.NewUnAuthorized("Invalid username/password combination")
		}
		return nil, utils.StandardInternalServerError()
	}

	// Check Password
	err = bcrypt.CompareHashAndPassword([]byte(foundUser.Password), []byte(req.Password))
	if err != nil {
		return nil, utils.NewUnAuthorized("Invalid username/password combination")
	}

	// Check verified email
	//if foundUser.EmailVerified == 0 {
	//	return nil, utils.NewBadRequest("Account not active, please check your inbox for a verification email", "EMAIL_NOT_ACTIVE")
	//}

	token, err := utils.CreateAuthToken(foundUser.ID)
	if err != nil {
		return nil, utils.StandardInternalServerError()
	}

	return token, nil
}

func (i *userService) Register(req *user.ReqRegister) utils.RestErrI {
	inUse, err := user.UserDao.EmailInUse(req.Email)

	if err != nil {
		return utils.StandardInternalServerError()
	}

	if inUse {
		return utils.NewBadRequest("This email is already associated to an account", "EMAIL_IN_USE")
	}

	// Hash Password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), 11)
	if err != nil {
		return utils.StandardInternalServerError()
	}

	// Generate Email Token
	EmailToken := uuid.New().String()

	_, err = user.UserDao.Create(&user.User{
		Name:       req.Name,
		Email:      req.Email,
		EmailToken: EmailToken,
		Password:   string(hashedPassword),
	})
	go user.UserDao.VerifyEmail(req.Email)

	if err != nil {
		return utils.StandardInternalServerError()
	}

	//go clients.SendMail("register.html", req.Email, "Welcome Email", struct {
	//	Token string
	//}{EmailToken})

	return nil
}
