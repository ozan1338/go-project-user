package userhandler

import (
	"net/http"
	"project/domain/user"
	"project/helpers"
	"project/service"
	"time"

	"project/pkg/jwt"

	"github.com/gorilla/mux"
)

var accesTokenDuration time.Duration = 15 * time.Minute

type userHandler struct {
	userService service.UserServiceInterface
	helpers helpers.HelpersInterface
	JWT jwt.Maker
}

func NewUserHandler(userService service.UserServiceInterface, helpers helpers.HelpersInterface,JWT jwt.Maker) *userHandler {
	return &userHandler{
		userService: userService,
		helpers: helpers,
		JWT: JWT,
	}
}

func (h userHandler) GetAllUser(w http.ResponseWriter, r *http.Request) {
	result, err := h.userService.GetAllUser()
	if err != nil {
		h.helpers.WriteResponse(w, err.GetStatus(), err)
		return
	}

	h.helpers.WriteResponse(w, http.StatusOK, result)
}

func (h userHandler) GetUserByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["user_id"]

	user, err := h.userService.GetByID(id)

	if err != nil {
		h.helpers.WriteResponse(w,err.GetStatus(),err)
		return
	}

	h.helpers.WriteResponse(w, http.StatusOK, user)
}

func (h userHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var user user.UsersRequest
	if err := h.helpers.ReadJSON(w,r,&user); err != nil {
		h.helpers.WriteResponse(w,err.GetStatus(),err)
		return
	}
	
	result, err := h.userService.CreateUser(user)
	if err != nil {
		h.helpers.WriteResponse(w,err.GetStatus(),err)
		return
	}

	

	jwtToken, _, err :=h.JWT.CreateToken(result.ID,accesTokenDuration)

	if err != nil {
		h.helpers.WriteResponse(w,err.GetStatus(),err)
		return
	}

	result.JWT = jwtToken

	h.helpers.WriteResponse(w,http.StatusOK, result)

}

func (h userHandler) LoginUser(w http.ResponseWriter, r *http.Request) {
	var user user.UsersRequest
	if err := h.helpers.ReadJSON(w,r,&user); err != nil {
		h.helpers.WriteResponse(w,err.GetStatus(),err)
		return
	}

	result, err := h.userService.LoginUser(user)
	if err != nil {
		h.helpers.WriteResponse(w,err.GetStatus(),err)
		return
	}

	jwtToken, _, err := h.JWT.CreateToken(result.ID, accesTokenDuration)
	if err != nil {
		h.helpers.WriteResponse(w,err.GetStatus(),err)
		return
	}

	result.JWT = jwtToken

	h.helpers.WriteResponse(w,http.StatusOK, result)
}