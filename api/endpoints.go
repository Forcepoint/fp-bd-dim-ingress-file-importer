package api

import (
	"encoding/json"
	"fmt"
	"fp-dim-aws-guard-duty-ingress/internal/config"
	"fp-dim-aws-guard-duty-ingress/internal/filehandler"
	"fp-dim-aws-guard-duty-ingress/internal/structs"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
	"strconv"
)

type HttpResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

func FileUpload(handler filehandler.FileHandler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		filename, err := filehandler.HandleFileUpload(r, handler)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(&HttpResponse{
				Status:  http.StatusBadRequest,
				Message: err.Error(),
			})
			return
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(&HttpResponse{
			Status:  http.StatusOK,
			Message: fmt.Sprintf("%s has been uploaded successfully.", filename),
		})
		return
	})
}

func SendHealth(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func ConfigEndpoint(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(config.GetConfig())
	} else if r.Method == "POST" {
		w.WriteHeader(http.StatusOK)
		configObj := &structs.ConfigObject{}
		json.NewDecoder(r.Body).Decode(configObj)
		time, err := strconv.Atoi(configObj.Values.UpdaterSchedule)
		if err != nil {
			logrus.Info(err)
			return
		} else {
			config.WriteConfig(time)
		}
		flags := os.O_CREATE | os.O_TRUNC | os.O_WRONLY
		filehandler.WriteFile(config.UrlsFile, flags, configObj.Values.Lists)
		return
	}
}
