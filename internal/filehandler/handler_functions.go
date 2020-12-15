package filehandler

import (
	"net/http"
)

func HandleFileUpload(r *http.Request, listHandler FileHandler) (string, error) {
	err := r.ParseMultipartForm(32 << 20)
	if err != nil {
		return "", err
	}
	file, handler, err := r.FormFile("file") //retrieve the file from form data
	if err != nil {
		return "", err
	}
	go listHandler.HandleUploadedFile(file)
	return handler.Filename, nil
}
