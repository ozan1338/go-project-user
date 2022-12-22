package user

import resError "project/util/errors_response"

type UsersRequest struct {
	Name string `json:"name,omitempty"`
	Email string `json:"email"`
	Password string `json:"password"`
}


type UsersResponse struct {
	ID int `json:"id"`
	Email string `json:"email"`
	Name string `json:"name"`
	JWT string `json:"jwt"`
}

func (u Users) ToDto() UsersResponse {
	return UsersResponse{
		u.ID,
		u.Email,
		u.Name,
		"",
	}
}

func (r UsersRequest) Validate(register bool) (resError.RespError) {
	if register {
		if r.Name == "" || r.Email == "" || r.Password == "" {
			return resError.NewBadRequestError("Please Input Name, Email and Password")
		}
	} else {
		if r.Email == "" || r.Password == "" {
			return resError.NewBadRequestError("Please Input Email and Password")
		}
	}

	return nil
}
