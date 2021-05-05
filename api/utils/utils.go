package utils

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
	"math"
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

type RestErr struct {
	Err string `json:"error"`
	StatusCode int `json:"status_code"`
	Type string `json:"type,omitempty"`
}

func (i *RestErr) Error() string {
	return i.Err
}

func (i *RestErr) Code() int {
	return i.StatusCode
}

func (i *RestErr) ErrType() string {
	return i.Type
}

// NewInternalServerError returns RestErrI with a status code of 500
func NewInternalServerError(message string, types ...string) RestErrI {
	return &RestErr{
		Err:        message,
		StatusCode: http.StatusInternalServerError,
		Type:       strings.Join(types, ", "),
	}
}

// NewBadRequest returns RestErrI with a status code of 400
func NewBadRequest(message string, types ...string) RestErrI {
	return &RestErr{
		Err:        message,
		StatusCode: http.StatusBadRequest,
		Type:       strings.Join(types, ", "),
	}
}

// NewUnprocessableEntity returns RestErrI with a status code of 422
func NewUnprocessableEntity(message string, types ...string) RestErrI {
	return &RestErr{
		Err:        message,
		StatusCode: http.StatusUnprocessableEntity,
		Type:       strings.Join(types, ", "),
	}
}

// NewUnAuthorized returns RestErrI with a status code of 401
func NewUnAuthorized(message string, types ...string) RestErrI {
	return &RestErr{
		Err:        message,
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

	config.OutputPaths = []string{logPath}

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
			e := RestErr.NewError(http.StatusUnprocessableEntity, "Invalid JSON body")
			c.JSON(e.StatusCode, e)
			return false
		}
		return true
	}
}

func GinShouldBindFormData(data interface{}) func(c *gin.Context) bool {
	return func(c *gin.Context) bool {
		if err := c.Bind(data); err != nil {
			e := RestErr.NewError(http.StatusUnprocessableEntity, "Invalid form data")
			c.JSON(e.StatusCode, e)
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
			ExpiresAt: timeNow.Unix() + 3600,
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



func ErrorLogger(e error) {
	if e == nil {
		return
	}
	fmt.Printf("Timestamp: %v \nError: %v", time.Now(), e.Error())
	f, err := os.OpenFile("error_log.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return
	}
	str := fmt.Sprintf("Timestamp: %v \nError: %v \n", time.Now(), e.Error())
	_, _ = f.WriteString(str)
	f.Close()
}

func CalculateReadTime(text string) *int {
	p := bluemonday.StrictPolicy()
	text = p.Sanitize(text)
	// Average person reads 200-250 words per minute
	readTime := int(math.Ceil(float64(len(text)) / 200))
	return &readTime
}