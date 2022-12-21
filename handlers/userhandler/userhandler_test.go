package userhandler

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"project/domain/user"
	"project/helpers"
	mockService "project/mocks/service"
	resError "project/util/errors_response"

	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
)

var (
	s *mockService.MockUserServiceInterface
	h *userHandler
	// router *mux.Router
)

func setup(t *testing.T) func() {
	ctrl := gomock.NewController(t)
	s = mockService.NewMockUserServiceInterface(ctrl)
	helper := helpers.NewHelper()
	h = NewUserHandler(s,helper)
	// router = mux.NewRouter()
	return func ()  {
		// router = nil
		defer ctrl.Finish()
	}
}

func TestUserHandler(t *testing.T) {
	//arrange

	var teardown = setup(t)

	defer teardown()

	var users = []user.UsersResponse{
		{ID: 1,Email: "test@mail.com",Name: "ozan"},
	}

	test := []struct{
		name string
		method string
		json string
		paramID string
		handler http.HandlerFunc
		expectedStatus int
		stub func() *gomock.Call
	}{
		{"all user ok", http.MethodGet,``, ``, h.GetAllUser, http.StatusOK, func() *gomock.Call {
			return s.EXPECT().GetAllUser().Return(users,nil)
		}},
		{"all user error", http.MethodGet,``, "", h.GetAllUser, http.StatusBadRequest, func() *gomock.Call {
			return s.EXPECT().GetAllUser().Return(nil,resError.NewBadRequestError("database error"))
		}},
		{"get user ok", http.MethodGet,``, "1", h.GetUserByID, http.StatusOK, func() *gomock.Call {
			return s.EXPECT().GetByID("1").Return(&users[0],nil)
		}},
		{"get user error", http.MethodGet,``, "1", h.GetUserByID, http.StatusBadRequest, func() *gomock.Call {
			return s.EXPECT().GetByID("1").Return(nil,resError.NewBadRequestError("databae error"))
		}},
	}

	//act
	for _, item := range test{

		item.stub()

		var req *http.Request
		if item.json == "" {
			req,_ = http.NewRequest(item.method,"/", nil)
		} else {
			req, _ = http.NewRequest(item.method, "/", strings.NewReader(item.json))
		}

		if item.paramID != "" {
			t.Log(item.paramID)
			var val = map[string]string{
				"user_id":item.paramID,
			}
			req = mux.SetURLVars(req, val)
		}

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(item.handler)

		handler.ServeHTTP(rr, req)

		//assert
		if rr.Code  != item.expectedStatus {
			t.Errorf("%s: wrong status return; expected %d but got %d", item.name, item.expectedStatus, rr.Code)
		}
		
	}
}

