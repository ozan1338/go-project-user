package user

import (
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
		return resError.NewBadRequestError("error from server")
	}

	u.Password = string(hashedPass)

	return nil
}

func (u Users) CheckPassword(password string) bool {
	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password)); err != nil {
		return false
	}
	return true
}