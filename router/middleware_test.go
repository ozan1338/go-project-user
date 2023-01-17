package router

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"
)

var expiredToken = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhZG1pbiI6dHJ1ZSwiYXVkIjoiZXhhbXBsZS5jb20iLCJleHAiOjE2NzAwNDU5NDYsImlzcyI6ImV4YW1wbGUuY29tIiwibmFtZSI6IkpvaG4gRG9lIiwic3ViIjoiMSJ9.pLlTFqTsvUQ2xnSVyTsfp-Y2_qPyeLXBUtSSXAwXkkg"

func TestAuthHeader(t *testing.T) {
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})

	idTestUser := 1

	testToken , _, _ := jwtMaker.CreateToken(idTestUser,15 * time.Minute)

	var test = []struct{
		name string
		token string
		expectAuth bool
		setHeader bool
	} {
		{name: "valid", token: fmt.Sprintf("Bearer %s", testToken), expectAuth: true, setHeader: true},
		{name: "expired", token: fmt.Sprintf("Bearer %s", expiredToken), expectAuth: false, setHeader: true},
		{name: "not token", token: "", expectAuth: false, setHeader: false},
		{name: "not valid", token: fmt.Sprintf("Bearer %s test", testToken), expectAuth: false, setHeader: true},
		{name: "not valid", token: fmt.Sprintf("Test %s", testToken), expectAuth: false, setHeader: true},
	}

	for _, item := range test {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		
		if item.setHeader {
			req.Header.Set(authorization, item.token)
		}

		rr := httptest.NewRecorder()

		handlerToTest := authRequired(nextHandler)

		handlerToTest.ServeHTTP(rr, req)

		if item.expectAuth && rr.Code == http.StatusUnauthorized {
			t.Errorf("%s: got code 401 and should not have", item.name)
		}

		if !item.expectAuth && rr.Code != http.StatusUnauthorized {
			t.Errorf("%s: did not code 401 and shoul have", item.name)
		}
	}
}

func TestUploadFile(t *testing.T) {
	//arrange
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	uploadDir = "../upload/test-image"

	var test = []struct{
		name string
		filePath string
		expectedStatus int
	} {
		{"no error", "../upload/06.jpeg", http.StatusOK},
	}

	//act
	for _, item := range test {
		//specified file name for the form
		fieldName := "file"

		//create a bytes.Buffer to act as the request body
		body := new(bytes.Buffer)

		//craete new writer
		mw := multipart.NewWriter(body)

		file, err := os.Open(item.filePath)
		if err != nil {
			t.Fatal(err)
		}

		defer file.Close()

		w, err := mw.CreateFormFile(fieldName, item.filePath)
		if err != nil {
			t.Fatal(err)
		}

		if _,err := io.Copy(w, file); err != nil {
			t.Fatal(err)
		}

		mw.Close()

		req := httptest.NewRequest(http.MethodPost, "/upload", body)
		
		req.Header.Add("Content-Type", mw.FormDataContentType())

		rr := httptest.NewRecorder()

		handlerToTest := uploadFiles(nextHandler)

		handlerToTest.ServeHTTP(rr, req)

		if rr.Code != item.expectedStatus {
			t.Errorf("%s: wrong status code; expected %d but got %d", item.name, item.expectedStatus, rr.Code)
		}

		_ = os.Remove("../upload/test-image/06.jpeg")
	}
}