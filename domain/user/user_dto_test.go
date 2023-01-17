package user

import (
	"reflect"
	"testing"
)

func TestToDto(t *testing.T) {
	//arrange and act
	testLoginResponse := UsersResponse{ID: 1, Email: "example.com", Name: "ozan", JWT: ""}
	loginResponse := user.ToDto()

	if match := reflect.DeepEqual(testLoginResponse, loginResponse); !match {
		t.Errorf("return of func ToDto Is not expected")
	}
}

func TestValidateReqRegister(t *testing.T) {
	//arrange
	testLoginRequest := UsersRequest{
		Name: "test",
		Email: "test@example.com",
		Password: "123",
	}

	//act
	err := testLoginRequest.Validate(true)

	//assert
	if err != nil {
		t.Errorf("expected no error but got one %s", err.GetMessage())
	}

	testLoginRequest.Name = ""

	err = testLoginRequest.Validate(true)

	if err == nil {
		t.Errorf("expected err but got nothing")
	}
}

func TestValidateReqLogin(t *testing.T) {
	//arrange
	testLogingRequest := UsersRequest{
		Name: "",
		Email: "test@example.com",
		Password: "123",
	}

	//act
	err := testLogingRequest.Validate(false)

	//assert
	if err != nil {
		t.Errorf("Expected no error but got one %s",err.GetMessage())
	}

	testLogingRequest.Password = ""

	err = testLogingRequest.Validate(true)

	if err == nil {
		t.Errorf("expected err but got nothing")
	}
}