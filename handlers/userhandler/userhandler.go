package userhandler

import (
	"net/http"
	"project/helpers"
	"project/service"

	"github.com/gorilla/mux"
)

type userHandler struct {
	userService service.UserServiceInterface
	helpers helpers.HelpersInterface
}

func NewUserHandler(userService service.UserServiceInterface, helpers helpers.HelpersInterface) *userHandler {
	return &userHandler{
		userService: userService,
		helpers: helpers,
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