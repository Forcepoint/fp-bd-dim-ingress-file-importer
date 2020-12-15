package main

import (
	"fmt"
	"fp-dim-aws-guard-duty-ingress/internal/filehandler"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"strings"
	"testing"
)

func TestDiffFiles(t *testing.T) {
	differ := filehandler.NewFileDiffer()

	data, err := ioutil.ReadFile("./testlists/blocklist-a.txt")

	if err != nil {
		t.Error(err)
	}

	// Local blocklist goes first, then the file to be compared against it, the difference is returned
	diff := differ.DiffFiles("./testlists/blocklist-b.txt", strings.Split(string(data), "\n"))

	fmt.Println(diff)
	lif := len(diff)
	assert.False(t, lif == 0)
	assert.True(t, lif == 3)
	assert.Equal(t, []string{"app-c.net", "zvooq.eu", "zvuker.net"}, diff)

	data, err = ioutil.ReadFile("./testlists/blocklist-a.txt")

	if err != nil {
		t.Error(err)
	}

	diff = differ.DiffFiles("./lists/blocklist.txt", strings.Split(string(data), "\n"))

	fmt.Println(diff)
	lif = len(diff)
	assert.False(t, lif == 0)
	assert.True(t, lif == 1)
	assert.Equal(t, []string{"app-c.net"}, diff)
}
