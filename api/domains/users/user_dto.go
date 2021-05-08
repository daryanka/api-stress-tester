package users

type ReqLogin struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type ReqRegister struct {
	Email                string `json:"email" validate:"required,email,max=254"`
	Name                 string `json:"name" validate:"required,max=254"`
	Password             string `json:"password" validate:"required,max=254,min=6"`
	PasswordConfirmation string `json:"password_confirmation" validate:"required,max=254,min=6"`
}

type User struct {
	ID            int64  `json:"id" db:"id"`
	Name          string `json:"name" db:"name"`
	Email         string `json:"email" db:"email"`
	EmailToken    string `json:"-" db:"email_token"`
	Password      string `json:"-" db:"password"`
	EmailVerified int    `json:"email_verified" db:"email_verified"`
}
