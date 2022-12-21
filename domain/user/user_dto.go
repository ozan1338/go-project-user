package user

import resError "project/util/errors_response"

type Users struct {
	ID       int    `db:"id"`
	Name     string `db:"name"`
	Email    string `db:"email"`
	Password string `db:"password"`
}

type UsersLoginRequest struct {
	Email string `json:"email"`
	Password string `json:"password"`
}

type UsersResponse struct {
	ID int `json:"id"`
	Email string `json:"email"`
	Name string `json:"name"`
}

func (u Users) ToDto() UsersResponse {
	return UsersResponse{
		u.ID,
		u.Email,
		u.Name,
	}
}

func (r UsersLoginRequest) Validate() (bool, resError.RespError) {
	if r.Email == "" || r.Password == "" {
		return false,resError.NewBadRequestError("Please Input Email and Password")
	}

	return true, nil
}
