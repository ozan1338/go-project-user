package service

import (
	"project/domain/user"
	"project/repo"
	resError "project/util/errors_response"
)

//go:generate mockgen -destination=../mocks/service/mockUserService.go -package=service project/service UserServiceInterface
type UserServiceInterface interface {
	GetAllUser() ([]user.UsersResponse, resError.RespError)
	GetByID(string) (*user.UsersResponse, resError.RespError)
	CreateUser(user.UsersRequest) (*user.UsersResponse, resError.RespError)
}

type UserService struct {
	repo repo.UserRepo
}

func NewUserService(userRepo repo.UserRepo) UserServiceInterface {
	return &UserService{repo: userRepo}
}

func (s UserService) CreateUser(u user.UsersRequest) (*user.UsersResponse, resError.RespError) {

	if err := u.Validate(true); err != nil {
		return nil, err
	}

	var user user.Users
	user.Email = u.Email
	user.Name = u.Name
	user.Password = u.Password
	if err := user.HashPassword(); err != nil {
		return nil, err
	}

	id, err := s.repo.CreateUser(user);
	if err != nil {
		return nil, err
	}
	
	user.ID = id

	result := user.ToDto()

	return &result, nil
}

func (s UserService) GetAllUser() ([]user.UsersResponse, resError.RespError) {
	result, err :=s.repo.GetAllUser()
	if err != nil {
		return nil, err
	}

	var response []user.UsersResponse

	for _, u := range result{
		response = append(response, u.ToDto())
	}

	return response,err
}

func (s UserService) GetByID(id string) (*user.UsersResponse, resError.RespError) {
	u, err := s.repo.GetUserByID(id)

	if err != nil {
		return nil, err
	}

	response := u.ToDto()

	return &response,nil
}