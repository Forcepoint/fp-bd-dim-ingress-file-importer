package main

import (
	"fmt"
	"fp-dim-aws-guard-duty-ingress/api"
	"fp-dim-aws-guard-duty-ingress/internal"
	"fp-dim-aws-guard-duty-ingress/internal/config"
	"fp-dim-aws-guard-duty-ingress/internal/filehandler"
	"fp-dim-aws-guard-duty-ingress/internal/hooks"
	"github.com/go-co-op/gocron"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"net/http"
	"os"
	"time"
)

func main() {
	InitLogrus()
	config.InitFiles()
	config.InitConfig()
	internal.Register()

	router := mux.NewRouter()

	differ := filehandler.NewFileDiffer()

	handler := filehandler.NewListFileHandler(
		"/lists/",
		"blocklist.txt",
		"tempfile.txt",
		"urls.txt",
		differ)

	router.Handle("/run", api.FileUpload(handler)).Methods("POST", "OPTIONS")
	router.HandleFunc("/health", api.SendHealth).Methods("GET", "OPTIONS")
	router.HandleFunc("/config", api.ConfigEndpoint).Methods("GET", "POST", "OPTIONS")
	router.HandleFunc("/icon", func(res http.ResponseWriter, req *http.Request) {
		http.ServeFile(res, req, "./icon/fp.jpeg")
	})

	// defines a new scheduler that schedules and runs jobs
	s1 := gocron.NewScheduler(time.UTC)

	updateSchedule := viper.GetUint64("update-schedule")

	if updateSchedule == 0 {
		updateSchedule = 12
	}

	s1.Every(updateSchedule).Hours().StartImmediately().Do(handler.RunPeriodicUpdater)

	// scheduler starts running jobs and current thread continues to execute
	s1.StartAsync()

	logrus.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", os.Getenv("LOCAL_PORT")), router))
}

func InitLogrus() {
	logrus.SetLevel(logrus.InfoLevel)
	// Show where error was logged, function, line number, etc.
	logrus.SetReportCaller(true)

	// Output to stdout and logfile
	logrus.SetOutput(os.Stdout)

	logrus.AddHook(&hooks.LoggingHook{})
}
