package user

import (
	"errors"
	"github.com/go-playground/validator/v10"
	"github.com/uptrace/bun"
	"golang.org/x/crypto/bcrypt"
	"log"
	"regexp"
	"unicode"
)

type User struct {
	bun.BaseModel `bun:"table:users"`
	ID            int64  `bun:",pk,autoincrement" json:"id"`
	Name          string `bun:",notnull" json:"name" validate:"name"`
	Email         string `bun:",unique,notnull" json:"email" validate:"email"`
	Password      string `bun:",notnull" json:"password" validate:"password"`
}

type LoginUser struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (u *User) HashPassword() {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal("error hashing password")
	}
	u.Password = string(hashedPassword)
}

func (l *LoginUser) CheckPassword(hashPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashPassword), []byte(l.Password))
}

func NameValidator(fl validator.FieldLevel) bool {
	str, ok := fl.Field().Interface().(string)
	return ok && str != ""
}

func EmailValidator(fl validator.FieldLevel) bool {
	email, ok := fl.Field().Interface().(string)
	if !ok {
		return false
	}
	re := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return re.MatchString(email)
}

// PasswordValidator enforces: minimum 8 characters, at least one digit, at least one uppercase letter
func PasswordValidator(fl validator.FieldLevel) bool {
	password, ok := fl.Field().Interface().(string)
	if !ok {
		return false
	}
	if len(password) < 8 {
		return false
	}
	hasDigit := false
	hasUpper := false
	for _, c := range password {
		if unicode.IsDigit(c) {
			hasDigit = true
		}
		if unicode.IsUpper(c) {
			hasUpper = true
		}
	}
	return hasDigit && hasUpper
}

func UserValidationError(err error) map[string]string {
	var validationErrors validator.ValidationErrors
	if !errors.As(err, &validationErrors) {
		return map[string]string{"error": "Unknown error"}
	}
	errorsMap := make(map[string]string)
	for _, e := range validationErrors {
		field := e.Field()
		switch e.Tag() {
		case "name":
			errorsMap[field] = "Provide your full name"
		case "email":
			errorsMap[field] = "Provide valid email address"
		case "password":
			errorsMap[field] = "Password must be at least 8 characters with one uppercase letter and one number"
		default:
			errorsMap[field] = "Invalid"
		}
	}
	return errorsMap
}