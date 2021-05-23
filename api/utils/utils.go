package utils

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
	"math/rand"
	"net/http"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"
)

var Logger *zap.SugaredLogger

type RestErrI interface {
	Error() string
	Code() int
	ErrType() string
}

type restErr struct {
	Err        string `json:"error"`
	StatusCode int    `json:"status_code"`
	Type       string `json:"type,omitempty"`
}

func (i *restErr) Error() string {
	return i.Err
}

func (i *restErr) Code() int {
	return i.StatusCode
}

func (i *restErr) ErrType() string {
	return i.Type
}

// NewInternalServerError returns RestErrI with a status code of 500
func NewInternalServerError(message string, types ...string) RestErrI {
	return &restErr{
		Err:        message,
		StatusCode: http.StatusInternalServerError,
		Type:       strings.Join(types, ", "),
	}
}

func StandardInternalServerError(types ...string) RestErrI {
	return &restErr{
		Err:        "Something went wrong, please try again later",
		StatusCode: http.StatusInternalServerError,
		Type: strings.Join(types, ", "),
	}
}

// NewBadRequest returns RestErrI with a status code of 400
func NewBadRequest(message string, types ...string) RestErrI {
	return &restErr{
		Err:        message,
		StatusCode: http.StatusBadRequest,
		Type:       strings.Join(types, ", "),
	}
}

// NewUnprocessableEntity returns RestErrI with a status code of 422
func NewUnprocessableEntity(message string, types ...string) RestErrI {
	return &restErr{
		Err:        message,
		StatusCode: http.StatusUnprocessableEntity,
		Type:       strings.Join(types, ", "),
	}
}

// NewUnAuthorized returns RestErrI with a status code of 401
func NewUnAuthorized(message string, types ...string) RestErrI {
	return &restErr{
		Err:        message,
		StatusCode: http.StatusUnauthorized,
		Type:       strings.Join(types, ", "),
	}
}

// StandardUnauthorized returns RestErrI with a status code of 401 and message Unauthorized
func StandardUnauthorized(types ...string) RestErrI {
	return &restErr{
		Err:        "Unauthorized",
		StatusCode: http.StatusUnauthorized,
		Type:       strings.Join(types, ", "),
	}
}

func InitLogger() {
	logPath := os.Getenv("LOG_PATH")
	mode := os.Getenv("MODE")
	config := zap.Config{}

	if mode == "PROD" {
		config = zap.NewProductionConfig()
	} else {
		config = zap.NewDevelopmentConfig()
	}

	config.OutputPaths = []string{logPath, "stdout"}

	logger, err := config.Build()
	if err != nil {
		panic(err.Error())
	}
	Logger = logger.Sugar()
}

type TokenWithClaims struct {
	Token     string `json:"token"`
	Expires   int64  `json:"expires"`
	ExpiresIn int    `json:"expires_in"`
}

func GinShouldPassAll(c *gin.Context, functions ...func(c *gin.Context) bool) bool {
	for _, el := range functions {
		if ok := el(c); !ok {
			return false
		}
	}
	return true
}

func GinShouldValidate(data interface{}) func(c *gin.Context) bool {
	return func(c *gin.Context) bool {
		if err := ValidateStruct(data); err != nil {
			c.JSON(http.StatusUnprocessableEntity, err)
			return false
		}
		return true
	}
}

func GinShouldBindJSON(data interface{}) func(c *gin.Context) bool {
	return func(c *gin.Context) bool {
		if err := c.ShouldBindJSON(&data); err != nil {
			e := NewUnprocessableEntity("There was an issue making the request, please try again later.", "INVALID_JSON")
			c.JSON(e.Code(), e)
			return false
		}
		return true
	}
}

func GinShouldBindFormData(data interface{}) func(c *gin.Context) bool {
	return func(c *gin.Context) bool {
		if err := c.Bind(data); err != nil {
			e := NewUnprocessableEntity("There was an issue making the request, please try again later.", "INVALID_FORM_DATA")
			c.JSON(e.Code(), e)
			return false
		}
		return true
	}
}

func GinGetPagination(c *gin.Context) (page int, perPage int) {
	p := c.DefaultQuery("page", "1")
	perP := c.DefaultQuery("per_page", "15")

	page, err := strconv.Atoi(p)
	if err != nil {
		page = 1
	}
	perPage, err = strconv.Atoi(perP)
	if err != nil || perPage > 100 {
		perPage = 15
	}
	return
}

func ValidateStruct(obj interface{}) map[string]string {
	v := validator.New()
	v.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]

		if name == "-" {
			return ""
		}

		return name
	})

	err := v.Struct(obj)

	if err == nil {
		return nil
	}

	errMap := make(map[string]string)
	for _, e := range err.(validator.ValidationErrors) {
		field := e.Namespace()
		fieldSplit := strings.Split(field, ".")
		fieldSplit = fieldSplit[1:]
		field = strings.Join(fieldSplit, ".")

		switch e.ActualTag() {
		case "len":
			errMap[field] = fmt.Sprintf("%v must have a length of %v", field, e.Param())
		case "required":
			errMap[field] = fmt.Sprintf("%v is a required field", field)
		case "max":
			errMap[field] = fmt.Sprintf("%v must have a length less than %v", field, e.Param())
		case "min":
			errMap[field] = fmt.Sprintf("%v must have a length greater than %v", field, e.Param())
		case "email":
			errMap[field] = fmt.Sprintf("%v must be a valid email", field)
		default:
			errMap[field] = fmt.Sprintf("%v is invalid", field)
		}
	}

	return errMap
}

func CreateAuthToken(userID int64) (*TokenWithClaims, error) {
	type JWTClaims struct {
		ID int64 `json:"id"`
		jwt.StandardClaims
	}
	timeNow := time.Now()
	customClaims := JWTClaims{
		ID: userID,
		StandardClaims: jwt.StandardClaims{
			//ExpiresAt: timeNow.Unix() + 3600,
			ExpiresAt: timeNow.Unix() + 10000000,
			IssuedAt:  timeNow.Unix(),
			Issuer:    "journal_api",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, customClaims)
	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))

	return &TokenWithClaims{
		Token:     tokenString,
		Expires:   customClaims.StandardClaims.ExpiresAt,
		ExpiresIn: 3600,
	}, err
}

func RandStringRunes(n int) string {
	var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}
