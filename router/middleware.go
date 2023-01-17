package router

import (
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"project/helpers"
	"project/log"
	"project/pkg/jwt"
	"project/upload"
	resError "project/util/errors_response"
	"strings"
)

var (
	authorization = "Authorization"
	unauthorized = "unauthorized"
	jwtMaker = jwt.NewJWTMaker("secret")
	helper = helpers.NewHelper()
)

func authRequired(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//get the authorization header
		authHeader := r.Header.Get(authorization)

		//check
		if err := checkBearerToken(authHeader,r); err != nil {
			helper.WriteResponse(w,err.GetStatus(),err)
			return
		}
		next.ServeHTTP(w,r)
	})
}

func checkBearerToken(authHeader string, r *http.Request) resError.RespError {
	//sanity check

	if authHeader == "" {
		return resError.NewRespError("no auth header", http.StatusUnauthorized, unauthorized)
	}

	//split header on spaces
	headerSplit := strings.Split(authHeader, " ")
	if len(headerSplit) != 2 {
		return resError.NewRespError("invalid auth header", http.StatusUnauthorized, unauthorized)
	}

	//check to see if we have word "Bearer"
	if headerSplit[0] != "Bearer" {
		return resError.NewRespError("unauthorized: no Bearer", http.StatusUnauthorized, unauthorized)
	}

	token := headerSplit[1]

	_, err := jwtMaker.VerifyToken(token)
	if err != nil {
		return err
	}

	return nil
}

type UploadedFile struct {
	OriginalName 	string
	FileSize int64
}

// var uploadDir = "../upload"
var uploadDir = upload.GetFolderName()

func uploadFiles(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var uploadFiles []*UploadedFile
		var responseErr resError.RespError

		var imageSize int64 = 1024 * 1024 * 5

		err := r.ParseMultipartForm(imageSize)
		if err != nil {
			responseErr = resError.NewBadRequestError("file image is to big")
			helper.WriteResponse(w,responseErr.GetStatus(),responseErr)
			return
		}

		var infile multipart.File
		var outfile *os.File

		for _, fHeaders := range r.MultipartForm.File {
			for _, hdr := range fHeaders {
				uploadFiles , err = func(uploadedFiles []*UploadedFile) ([]*UploadedFile, error) {
					var uploadFile UploadedFile
					infile, err = hdr.Open()
					if err != nil {
						return nil, err
					}
					// defer infile.Close()
					

					uploadFile.OriginalName = hdr.Filename

					// var outfile *os.File
					// defer outfile.Close()

					buff := make([]byte, 512)
					_, err = infile.Read(buff)
					if err != nil {
						return nil,err
					}
					filetype := http.DetectContentType(buff)
					if filetype != "image/jpeg" && filetype != "image/png" && filetype != "image/jpg" && filetype != "application/octet-stream" {
						errMsg := fmt.Sprintf("format %s is not allowed, Please upload a JPEG/JPG or PNG", filetype)
						return nil, errors.New(errMsg)
					}

					outfile, err = os.Create(filepath.Join(uploadDir, uploadFile.OriginalName))

					if err != nil {
						return nil, err
					} else {
						fileSize, err := io.Copy(outfile, infile)
						if err != nil {
							return nil, err
						}

						uploadFile.FileSize = fileSize
					}

					uploadedFiles = append(uploadedFiles, &uploadFile)

					return uploadedFiles, nil
				} (uploadFiles)

				if err != nil {
					log.Error("error when create image", err)
				}
			}
			defer infile.Close()
			defer outfile.Close()
			if err != nil {
				responseErr = resError.NewBadRequestError("error when create image")
				helper.WriteResponse(w,responseErr.GetStatus(),responseErr)
				return
			}
		}

		next.ServeHTTP(w,r)
	})
}