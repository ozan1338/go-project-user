package router

import (
	"net/http"
	"testing"

	"github.com/gorilla/mux"
)

func TestUserRoutes(t *testing.T) {
	//arrange
	r := mux.NewRouter().StrictSlash(true)
	RouterInit(r.PathPrefix("/api/v1").Subrouter())

	var test = []struct{
		path string
		method string
	} {
		{"/api/v1/", http.MethodGet},
		{"/api/v1/get-all", http.MethodGet},
		{"/api/v1/get-user/{user_id:[0-9]+}", http.MethodGet},
		{"/api/v1/create-user", http.MethodPost},
		{"/api/v1/login", http.MethodGet},
	}

	for _, item := range test {
		registered, _ := registerRoute(item.path,item.method,r)
		if !registered {
			t.Errorf("route %s is not register", item.path)
		}
	}
}

func registerRoute(testPath string, testMethod string,r *mux.Router) (bool,string) {
	found := false
	isMethodWrong := ""

	r.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		path, _ := route.GetPathTemplate()
		if testPath == path {
			
			found = true
		}

		// method, _ := route.GetMethods()

		// if len(method) > 0 {
		// 	if strings.EqualFold(testMethod, method[0]) {
		// 		found = true
		// 	} else {
		// 		isMethodWrong = method[0]
		// 	}
		// }

		return nil
	})

	return found, isMethodWrong
}