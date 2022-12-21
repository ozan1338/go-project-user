package router

import (
	"net/http"

	"project/handlers/pinghandler"
	"project/handlers/userhandler"
	"project/helpers"
	"project/pkg/postgresql"
	"project/repo"
	"project/service"

	"github.com/gorilla/mux"
)

func userRouter(r *mux.Router) {
	//TODO
	//1.add repository db
	db := postgresql.Conn
	repo := repo.RepoUser(db)

	//2. add helper
	helper := helpers.NewHelper()

	//4. add service
	userService := service.NewUserService(repo)

	//3.add handlers for user
	h := pinghandler.PingHandler(helper)
	u := userhandler.NewUserHandler(userService, helper)

	//serve the route
	r.HandleFunc("/", h.Ping).Methods(http.MethodGet)
	r.HandleFunc("/get-all", u.GetAllUser).Methods(http.MethodGet)
	r.HandleFunc("/get-user/{user_id:[0-9]+}",u.GetUserByID).Methods(http.MethodGet)
}