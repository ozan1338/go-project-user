package userhandler

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"project/domain/user"
	"project/helpers"
	mockJwt "project/mocks/pkg/jwt"
	mockService "project/mocks/service"
	_ "project/pkg/jwt"
	resError "project/util/errors_response"

	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
)

var (
	s *mockService.MockUserServiceInterface
	h *userHandler
	jwtMaker *mockJwt.MockMaker
	// router *mux.Router
)

func setup(t *testing.T) func() {
	ctrl := gomock.NewController(t)
	s = mockService.NewMockUserServiceInterface(ctrl)
	helper := helpers.NewHelper()
	// jwtMaker := jwt.NewJWTMaker("secret")
	jwtMaker  = mockJwt.NewMockMaker(ctrl)
	h = NewUserHandler(s,helper,jwtMaker)
	// router = mux.NewRouter()
	return func ()  {
		// router = nil
		defer ctrl.Finish()
	}
}

var users = []user.UsersResponse{
	{ID: 1,Email: "test@mail.com",Name: "ozan", JWT: ""},
}

var usersRequest = user.UsersRequest{
	Name: "ozan",
	Email: "test@mail.com",
	Password: "123",
}

func TestUserHandlerLoginUser(t *testing.T) {
	//arrange 
	var teardown = setup(t)
	defer teardown()

	test := []struct{
		name string
		json string
		expectedStatus int
		stubService func() *gomock.Call
		stubJwt func() *gomock.Call
	} {
		{
			name:"login user ok",
			json: `{
				"email":"test@mail.com",
				"password":"123"
			}`,
			expectedStatus: http.StatusOK,
			stubService: func() *gomock.Call {
				var testUser = usersRequest
				testUser.Name = ""
				return s.EXPECT().LoginUser(testUser).Return(&users[0],nil)
			},
			stubJwt: func() *gomock.Call {
				return jwtMaker.EXPECT().CreateToken(users[0].ID, 15 * time.Minute)
			},
		},
		{
			name:"login user bad json",
			json: `{
				"email":"test@mail.com",
				"password":123
			}`,
			expectedStatus: http.StatusBadRequest,
			stubService: func() *gomock.Call {
				return nil
			},
			stubJwt: func() *gomock.Call {
				return nil
			},
		},
		{
			name:"login user error",
			json: `{
				"email":"test@mail.com",
				"password":"123"
			}`,
			expectedStatus: http.StatusBadRequest,
			stubService: func() *gomock.Call {
				var testUser = usersRequest
				testUser.Name = ""
				return s.EXPECT().LoginUser(testUser).Return(nil,resError.NewBadRequestError("some error"))
			},
			stubJwt: func() *gomock.Call {
				return nil
			},
		},
		{
			name:"login user jwt error",
			json: `{
				"email":"test@mail.com",
				"password":"123"
			}`,
			expectedStatus: http.StatusUnauthorized,
			stubService: func() *gomock.Call {
				var testUser = usersRequest
				testUser.Name = ""
				return s.EXPECT().LoginUser(testUser).Return(&users[0],nil)
			},
			stubJwt: func() *gomock.Call {
				return jwtMaker.EXPECT().CreateToken(users[0].ID, 15 * time.Minute).Return("",nil,resError.NewRespError("some error", http.StatusUnauthorized,"unauthorized"))
			},
		},
	}

	for _, item := range test {
		item.stubJwt()
		item.stubService()

		var req *http.Request
		req, _ = http.NewRequest(http.MethodPost,"/", strings.NewReader(item.json))

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(h.LoginUser)
		handler.ServeHTTP(rr,req)

		if rr.Code != item.expectedStatus {
			t.Errorf("%s : expected %d but got %d",item.name,item.expectedStatus,rr.Code)
		}
	}

}

func TestUserHandlerCreateUser(t *testing.T) {
	//arrange 
	var teardown = setup(t)
	defer teardown()

	test := []struct{
		name string
		json string
		expectedStatus int
		stubService func() *gomock.Call
		stubJwt func() *gomock.Call
	} {
		{"create user ok",
		`{
			"name":"ozan",
			"email":"test@mail.com",
			"password":"123"
		}`, http.StatusOK, 
		func() *gomock.Call { 
			return s.EXPECT().CreateUser(usersRequest).Return(&users[0],nil) 
		},
		func() *gomock.Call {
			return jwtMaker.EXPECT().CreateToken(users[0].ID,15 * time.Minute)
		}},
		{"create user bad json",
		`{
			"name":"ozan",
			"email":"test@mail.com",
			"password":123
		}`, http.StatusBadRequest, 
		func() *gomock.Call { 
			return nil
		},func() *gomock.Call {
			return nil
		}},
		{"create user error",
		`{
			"name":"ozan",
			"email":"test@mail.com",
			"password":"123"
		}`, http.StatusBadRequest, 
		func() *gomock.Call { 
			return s.EXPECT().CreateUser(usersRequest).Return(nil,resError.NewBadRequestError("some error"))
		},func() *gomock.Call {
			return nil
		}},
		{"create user jwt error",
		`{
			"name":"ozan",
			"email":"test@mail.com",
			"password":"123"
		}`, http.StatusUnauthorized, 
		func() *gomock.Call { 
			return s.EXPECT().CreateUser(usersRequest).Return(&users[0],nil)
		},func() *gomock.Call {
			return jwtMaker.EXPECT().CreateToken(users[0].ID,15 * time.Minute).Return("",nil,resError.NewRespError("some error", http.StatusUnauthorized, "jwt error"))
		}},
	}

	//act
	for _, item := range test{
		item.stubService()
		item.stubJwt()

		var req *http.Request
		req, _ = http.NewRequest(http.MethodPost,"/", strings.NewReader(item.json))

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(h.CreateUser)
		handler.ServeHTTP(rr,req)

		if rr.Code != item.expectedStatus {
			t.Errorf("%s: wrong status return; expected %d but got %d", item.name, item.expectedStatus, rr.Code)
		}
	}
}

func TestUserHandler(t *testing.T) {
	//arrange

	var teardown = setup(t)

	defer teardown()



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

