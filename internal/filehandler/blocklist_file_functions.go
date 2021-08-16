package filehandler

import (
	"bufio"
	"github.com/google/gonids"
	"fp-dim-aws-guard-duty-ingress/internal"
	"fp-dim-aws-guard-duty-ingress/internal/structs"
	"github.com/sirupsen/logrus"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"regexp"
	"strings"
	"sync"
)

const WindowsLineBreak = "\r\n"
const UnixLineBreak = "\n"

var ipRegex = regexp.MustCompile("^\\b\\d{1,3}\\.\\d{1,3}\\.\\d{1,3}\\.\\d{1,3}\\b$")
var urlRegex = regexp.MustCompile("^(http(s)?://)?(www\\.)?([-a-zA-Z0-9@:%_+~#=]{2,256}\\.[a-z]{2,256}\\b([-a-zA-Z0-9@:%_+~#?&/=]*))+(\\.[a-z]{2,6}\\b([-a-zA-Z0-9@:%_+~#?&/=]*))?$")
var domainRegex = regexp.MustCompile("^[a-zA-Z0-9][a-zA-Z0-9-]{1,61}[a-zA-Z0-9]\\.[a-zA-Z]{2,}$")
var snortRegexIgnoringRuleOptions = regexp.MustCompile("\\([^)]*\\)")

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
	for _, line := range strings.Split(strings.ReplaceAll(string(urlFile), WindowsLineBreak, UnixLineBreak), UnixLineBreak) {
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

	fileLines := sanitiseLines(strings.Split(strings.ReplaceAll(string(data), WindowsLineBreak, UnixLineBreak), UnixLineBreak))

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
	source := "Generic Data Importer"

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

		if isIP(element) {
			item.Type = "IP"
		} else if isDomain(element) {
			item.Type = "DOMAIN"
		} else if isUrl(element) {
			item.Type = "URL"
		} else if isSnort(element) {
			item.Type = "SNORT"
		} else {
			continue
		}

		item.Value = element
		item.Safe = false

		wrapper.Items = append(wrapper.Items, *item)
	}

	return internal.PushDataToController(wrapper)
}

func isIP(val string) bool {
	return ipRegex.MatchString(val)
}

func isDomain(val string) bool {
	return domainRegex.MatchString(val)
}

func isUrl(val string) bool {
	return urlRegex.MatchString(val)
}

func isSnort(val string) bool {
	_, err := gonids.ParseRule(snortRegexIgnoringRuleOptions.ReplaceAllString(val, "()"))
	if err != nil {
		return false
	}
	return true
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
