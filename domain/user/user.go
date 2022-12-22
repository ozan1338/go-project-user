package user

import (
	"net/http"
	resError "project/util/errors_response"

	"golang.org/x/crypto/bcrypt"
)

type Users struct {
	ID       int    `db:"id"`
	Name     string `db:"name"`
	Email    string `db:"email"`
	Password string `db:"password"`
}

func (u *Users) HashPassword() resError.RespError {
	hashedPass, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return resError.NewRespError("password not match", http.StatusUnauthorized,"bad credentials")
	}

	u.Password = string(hashedPass)

	return nil
}