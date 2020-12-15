package filehandler

import (
	"bufio"
	"fp-dim-aws-guard-duty-ingress/internal"
	"fp-dim-aws-guard-duty-ingress/internal/structs"
	"github.com/sirupsen/logrus"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
)

type FileHandler interface {
	RunPeriodicUpdater()
	HandleUploadedFile(multipart.File) error
}

type listFileHandler struct {
	mu                *sync.Mutex
	FileDirectory     string
	BlocklistFilename string
	TempFileName      string
	UrlListFilename   string
	differ            Differ
}

func NewListFileHandler(fileDirectory, blockListFile, tempFile, urlListFile string, differ Differ) *listFileHandler {
	return &listFileHandler{
		mu:                &sync.Mutex{},
		FileDirectory:     fileDirectory,
		BlocklistFilename: blockListFile,
		TempFileName:      tempFile,
		UrlListFilename:   urlListFile,
		differ:            differ,
	}
}

func (l *listFileHandler) RunPeriodicUpdater() {
	logrus.Info("Running updater...")
	urlFile, err := ioutil.ReadFile(l.FileDirectory + l.UrlListFilename)

	if err != nil {
		logrus.Error(err)
	}

	// Iterate over urls in url file
	for _, line := range strings.Split(string(urlFile), "\n") {
		if line == "" {
			continue
		}
		logrus.Info("Running for URL: " + line)
		// Pull the data from the URL and save to a temp file
		remoteLines := l.pullRemoteList(line)

		// Compare temp file to local blocklist
		diff := l.differ.DiffFiles(l.FileDirectory+l.BlocklistFilename, remoteLines)

		go func(diffLines []string) {
			err = writeToLocalFile(l.mu, l.FileDirectory+l.BlocklistFilename, diffLines)
			if err != nil {
				logrus.Error(err)
			}
		}(diff)

		_, _, err = ProcessTextFile(diff)

		if err != nil {
			logrus.Error(err)
			return
		}
	}
	return
}

func (l *listFileHandler) HandleUploadedFile(header multipart.File) error {
	logrus.Info("Handling uploaded file")
	data, err := ioutil.ReadAll(header)

	if err != nil {
		return err
	}

	fileLines := sanitiseLines(strings.Split(string(data), "\n"))

	// Compare temp file to local blocklist
	diff := l.differ.DiffFiles(l.FileDirectory+l.BlocklistFilename, fileLines)

	go func(diffLines []string) {
		err = writeToLocalFile(l.mu, l.FileDirectory+l.BlocklistFilename, diffLines)
		if err != nil {
			logrus.Error(err)
		}
	}(diff)

	_, _, err = ProcessTextFile(diff)

	if err != nil {
		logrus.Error(err)
		return err
	}

	return nil
}

func ProcessTextFile(elements []string) (*io.ReadCloser, int, error) {
	wrapper := structs.ProcessedItemsWrapper{}

	svcName := os.Getenv("MODULE_SVC_NAME")
	source := "Forcepoint File Importer"

	for _, element := range elements {
		if element == "" {
			continue
		}

		// ignore comments in files
		if strings.HasPrefix(element, "#") {
			continue
		}

		item := new(structs.ProcessedItem)
		item.ServiceName = svcName
		item.Source = source

		if IsIP(element) {
			item.Type = "IP"
		} else {
			val, err := isUrlOrDomain(element)
			if err != nil {
				logrus.Error(err)
				continue
			}
			item.Type = val
		}

		item.Value = element
		item.Safe = false

		wrapper.Items = append(wrapper.Items, *item)
	}

	return internal.PushDataToController(wrapper)
}

func IsIP(host string) bool {
	return net.ParseIP(host) != nil
}

func isUrlOrDomain(toTest string) (string, error) {
	_, err := url.ParseRequestURI(toTest)

	if err != nil {
		return "DOMAIN", nil
	}

	return "URL", nil
}

func writeToLocalFile(mu *sync.Mutex, filename string, dataToWrite []string) error {
	mu.Lock()
	defer mu.Unlock()
	f, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		logrus.Fatal(err)
	}
	defer f.Close()
	w := bufio.NewWriter(f)
	for _, s := range dataToWrite {
		_, err = w.WriteString(s + "\n")
	}

	return w.Flush()
}

func (l *listFileHandler) pullRemoteList(remoteUrl string) []string {
	resp, err := http.Get(remoteUrl)

	if err != nil {
		logrus.Error(err)
		return []string{}
	}

	data, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		logrus.Error(err)
		return []string{}
	}

	return sanitiseLines(strings.Split(string(data), "\n"))
}

func sanitiseLines(data []string) []string {
	for i, s := range data {
		data[i] = strings.TrimRight(s, ",")
	}
	return data
}
