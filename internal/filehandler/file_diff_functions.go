package filehandler

import (
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"strings"
)

type Differ interface {
	DiffFiles(string, []string) []string
}

type fileDiffer struct {
}

func NewFileDiffer() *fileDiffer {
	return &fileDiffer{}
}

func (d *fileDiffer) DiffFiles(localBlocklist string, data []string) []string {
	// Original source data
	local, err := ioutil.ReadFile(localBlocklist)

	if err != nil {
		logrus.Error(err)
	}

	return difference(data, strings.Split(string(local), "\n"))
}

// difference returns the elements in `a` that aren't in `b`.
func difference(a, b []string) []string {
	mb := make(map[string]int, len(b))
	for _, x := range b {
		mb[x] += 1
	}
	var diff []string
	for _, x := range a {
		if _, found := mb[x]; !found {
			mb[x] -= 1
			diff = append(diff, x)
		}
	}
	return diff
}
