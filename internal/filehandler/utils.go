package filehandler

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"strings"
)

func CreateFile(filename string) {
	f, err := os.Create(filename)
	if err != nil {
		logrus.Fatal(fmt.Sprintf("There was an error while creating the file. %s", filename))
		return
	}
	defer f.Close()
}

func OpenFile(filename string, flag int) *os.File {
	file, err := os.OpenFile(filename, flag, 0666)
	if err != nil {
		logrus.Fatal(fmt.Sprintf("There was an error while opening the file. %s", filename))
	}
	return file
}

func ReadFileContents(file *os.File) []string {
	data, err := ioutil.ReadAll(file)
	if err != nil {
		logrus.Fatal(err)
	}
	return strings.Split(string(data), "\n")
}

func WriteFile(filename string, flag int, data string) {
	f := OpenFile(filename, flag)

	data = strings.ReplaceAll(data, " ", "")

	data = strings.ReplaceAll(data, "\n", "")

	list := strings.Split(data, ",")

	for _, val := range list {
		if val != "" {
			f.WriteString(val + "\n")
		}
	}
	return
}
