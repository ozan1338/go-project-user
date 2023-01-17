package repo

import (
	"database/sql"
	"fmt"
	"project/domain/user"
	"project/log"
	resError "project/util/errors_response"
	"strings"
)

//go:generate mockgen -destination=../mocks/repo/mockUserRepo.go -package=repo project/repo UserRepo
type UserRepo interface {
	GetAllUser() ([]user.Users,resError.RespError)
	GetUserByID(string) (*user.Users,resError.RespError)
	CreateUser(user.Users) (int,resError.RespError)
	GetUserByEmail(*user.Users) (*user.Users,resError.RespError)
}

const (
	queryGetAllUser = `select id,name,email,password from public.User`
	queryGetUserById = `select id, name, email from public.User where id=$1`
	queryCreateUser = `insert into public.User(name,email,password) values($1,$2,$3) returning id;`
	queryGetUserByEmail = `select id, password, name from public.User where email = $1;`
	errNoRows = "no rows in result set"
)

func RepoUser(db *sql.DB) UserRepo {
	return &repo{db: db}
}

func (r *repo) GetUserByEmail(u *user.Users) (*user.Users,resError.RespError) {
	stmt, err := r.db.Prepare(queryGetUserByEmail)
	if err != nil {
		log.Error("error when trying to prepare get user by email", err)
		return nil,resError.NewBadRequestError("database error")
	}
	defer stmt.Close()

	if err := stmt.QueryRow(u.Email).Scan(&u.ID,&u.Password,&u.Name); err != nil {
		if strings.Contains(err.Error(),errNoRows) {
			return nil,resError.NewBadRequestError("user not found")
		}
		log.Error("error when trying scan get user by email", err)
		return nil,resError.NewBadRequestError("database error")
	}

	return u,nil
}

func (r *repo) CreateUser(u user.Users) (int, resError.RespError) {
	stmt, err := r.db.Prepare(queryCreateUser)
	if err != nil {
		log.Error("error when trying to prepare create user statement", err)
		return 0, resError.NewBadRequestError("database error")
	}
	defer stmt.Close()

	if err := stmt.QueryRow(u.Name,u.Email,u.Password).Scan(&u.ID); err != nil {
		log.Error("error when trying to scan create user", err)
		return 0, resError.NewBadRequestError("database error")
	}
	
	return u.ID, nil
}

func (r *repo) GetAllUser() ([]user.Users,resError.RespError) {
	// var users []user.Users
	stmt, err := r.db.Prepare(queryGetAllUser)
	if err != nil {
		log.Error("error when trying to prepare get all user statement",err)
		return nil, resError.NewBadRequestError("database error")
	}
	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		log.Error("error when trying to get all user", err)
		return nil, resError.NewBadRequestError("database error")
	}
	defer rows.Close()

	result := make([]user.Users, 0)
	for rows.Next() {
		var user user.Users
		if err := rows.Scan(&user.ID,&user.Name,&user.Email,&user.Password); err != nil {
			log.Error("error when trying to scan user", err)
			return nil, resError.NewBadRequestError("database error")
		}
		result = append(result, user)
	}

	if len(result) == 0 {
		return nil, resError.NewBadRequestError("no user found")
	}

	return result,nil
}

func (r *repo) GetUserByID(id string) (*user.Users, resError.RespError) {
	var user user.Users
	stmt, err := r.db.Prepare(queryGetUserById)
	if err != nil {
		log.Error("error when trying to prepare get user stmt", err)
		return nil,resError.NewBadRequestError("database error")
	}

	defer stmt.Close()

	result := stmt.QueryRow(id)
	if getErr := result.Scan(&user.ID,&user.Name, &user.Email); getErr != nil {
		if strings.Contains(getErr.Error(), errNoRows) {
			return nil,resError.NewBadRequestError(fmt.Sprintf("not found user with given id %s", id))
		}
		log.Error("error when trying to prepare get user", err)
		return nil, resError.NewBadRequestError("database error")
	}

	return &user,nil
}