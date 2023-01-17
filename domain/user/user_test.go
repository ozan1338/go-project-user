package user

import "testing"

var testUser =  Users{
	ID: 1,
	Name: "ozan",
	Email: "example.com",
	Password: "123",
}

var user = Users{
	ID: testUser.ID,
	Name: testUser.Name,
	Email: testUser.Email,
	Password: testUser.Password,
}

func TestHashPassword(t *testing.T) {

	//act
	err := user.HashPassword()

	//assert
	if err != nil {
		t.Error("got error while hashing password: ", err )
	}

	if user.Password == testUser.Password {
		t.Errorf("the hash password should not same with plain password")
	}
}

func TestCheckHashPassword(t *testing.T) {
	//act
	ok := user.CheckPassword(testUser.Password)

	// assert
	if !ok {
		t.Errorf("password should be same with originall password")
	}

	ok = user.CheckPassword("098")

	if ok {
		t.Errorf("password shoul be not match")
	}
}