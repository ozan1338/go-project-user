package service

import (
	"project/domain/user"
	mockRepo "project/mocks/repo"
	"testing"

	resError "project/util/errors_response"

	"github.com/golang/mock/gomock"
)

var (
	r *mockRepo.MockUserRepo
	service UserServiceInterface
)

func setup(t *testing.T) func() {
	ctrl := gomock.NewController(t)
	r = mockRepo.NewMockUserRepo(ctrl)
	service = NewUserService(r)
	return func ()  {
		service = nil
		defer ctrl.Finish()
	}
}

func TestUserServiceLoginUser(t *testing.T) {
	//Arrange
	teardown := setup(t)
	defer teardown()

	loginRequest := user.UsersRequest{
		Name: "",
		Email: "test@mail.com",
		Password: "123",
	}

	defaultUsers := user.Users{
		ID: 0,
		Name: loginRequest.Name,
		Email: loginRequest.Email,
		Password: "",
	}

	user := user.Users{
		ID: 0,
		Name: loginRequest.Email,
		Email: loginRequest.Email,
		Password: "$2a$10$s6a55hgBPxiCVBkHZgvG1eZfGszAheLBPQfR9f6VRRv8slBGI5J9K",
	}


	test := []struct{
		name string
		stub func() *gomock.Call
		expectedErr bool
	} {
		{
			name: "no error",
			stub: func() *gomock.Call {
				return r.EXPECT().GetUserByEmail(&defaultUsers).Return(&user,nil)
			},
			expectedErr: false,
		},
		{
			name: "error",
			stub: func() *gomock.Call {
				return r.EXPECT().GetUserByEmail(&defaultUsers).Return(nil,resError.NewBadRequestError("some database error"))
			},
			expectedErr: true,
		},
		{
			name: "error validate user",
			stub: func() *gomock.Call {
				return nil
			},
			expectedErr: true,
		},
		{
			name: "error password not match",
			stub: func() *gomock.Call {
				return r.EXPECT().GetUserByEmail(&defaultUsers).Return(&user,nil)
			},
			expectedErr: true,
		},
	}

	for _,item := range test {
		//Act
		item.stub()
		if item.name == "error password not match" {
			loginRequest.Password = "0"
		} 

		if item.name == "error validate user" {
			loginRequest.Email = ""
			loginRequest.Password = ""
		}

		_,err := service.LoginUser(loginRequest)

		if item.expectedErr && err == nil {
			t.Errorf("%s:expected error but got nothing", item.name)
		}
		
		if !item.expectedErr && err != nil {
			t.Errorf("%s:expected no error but got %s", item.name,err.GetMessage())
		}

		loginRequest.Password = "123"
		loginRequest.Email = "test@mail.com"
	}
}

func TestUserServiceGetAll(t *testing.T) {
	//Arrange
	teardown := setup(t)
	defer teardown()

	defaultUsers := []user.Users{
		{
			ID: 1,
			Name: "ozan",
			Email: "akhmadfauzan@gmail.com",
			Password: "1234",
		},
		
	}

	test := []struct{
		name string
		stub func() *gomock.Call
		expectedErr bool
	} {
		{
			name: "No Error",
			stub: func() *gomock.Call {
				return r.EXPECT().GetAllUser().Return(defaultUsers,nil)
			},
			expectedErr: false,
		},
		{
			name: "Error",
			stub: func() *gomock.Call {
				return r.EXPECT().GetAllUser().Return(nil,resError.NewBadRequestError("database error"))
			},
			expectedErr: true,
		},
	}

	for _, item := range test {
		//Act
		item.stub()
		users,err := service.GetAllUser()
	
		//Assert
		if item.expectedErr && err == nil {
			t.Errorf("%s:expected error but got nothing", item.name)
		}
		
		if !item.expectedErr && err != nil {
			t.Errorf("%s:expected no error but got %s", err.GetMessage(), item.name)
		}

		if len(users) > 0 {
			if defaultUsers[0].Email != users[0].Email {
				t.Errorf("%s:email not same", item.name)
			}
		}

	}
}

func TestUserServiceGetById(t *testing.T) {
	//arrange
	teardown := setup(t)
	defer teardown()

	var defaultUser = user.Users{
		ID: 1,
		Name: "ozan",
		Email: "here@mail.com",
	}


	test := []struct{
		name string
		stub func () *gomock.Call
		expctedErr bool
		testID string
	} {
		{
			name: "No Error",
			stub: func() *gomock.Call {
				return r.EXPECT().GetUserByID("1").Return(&defaultUser,nil)
			},
			expctedErr: false,
			testID: "1",
		},
		{
			name: "Error",
			stub: func() *gomock.Call {
				return r.EXPECT().GetUserByID("0").Return(nil,resError.NewBadRequestError("database error"))
			},
			expctedErr: true,
			testID: "0",
		},
	}

	for _, item := range test {
		//act
		item.stub()

		user,err := service.GetByID(item.testID)
		
		//assert
		if item.expctedErr && err == nil {
			t.Errorf("%s:expected error but got nothing",item.name)
		}

		if !item.expctedErr && err != nil {
			t.Errorf("%s:expected no error but got one %s", item.name,err.GetMessage())
		}

		if user != nil {
			if user.ID != defaultUser.ID {
				t.Errorf("%s:expected id %d but got id %d", item.name, defaultUser.ID, user.ID)
			}
		}
	}
}

func TestUserServiceCreateUser(t *testing.T) {
	//arrange
	teardown := setup(t)
	defer teardown()

	var defaultUser = user.UsersRequest{
		Email: "test@mail.com",
		Name: "ozan",
		Password: "123",
	}

	test := []struct{
		name string
		stub func() *gomock.Call
		expectedErr bool
	} {
		{"no error", func() *gomock.Call {return r.EXPECT().CreateUser(gomock.Any()).Return(1,nil)}, false},
		{"error", func() *gomock.Call {return r.EXPECT().CreateUser(gomock.Any()).Return(0,resError.NewBadRequestError("database error"))}, true},
	}

	for _, item := range test {
		//act
		item.stub()

		user,err := service.CreateUser(defaultUser)

		//act
		if err != nil && !item.expectedErr {
			t.Errorf("%s: not expected error but got one: %s",item.name,err.GetMessage())
		}

		if err == nil && item.expectedErr {
			t.Errorf("%s: expected error but got nothing", item.name)
		}

		if user != nil {
			if user.JWT != "" {
				t.Errorf("jwt field initially doesnt have value but got %s", user.JWT)
			}
		}
	}
}