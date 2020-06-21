package utils

import (
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"testing"
)

func TestEnsureFileExists(t *testing.T) {
	tmp := createTmpBranchProtectionConfigFile()

	exists, err := FileExists(tmp.Name())

	if err != nil {
		t.Fatal("not possible verify if file exists")
	}

	if !exists {
		t.Fatal("file not found")
	}
}

func TestEnsureFileNotExists(t *testing.T) {
	fileName := getRandomFileName()

	exists, err := FileExists(fileName)
	if err != nil {
		t.Fatal("not possible verify if file exists")
	}

	if exists {
		t.Fatal("expected not to find a file")
	}
}

func getRandomFileName() string {
	var max, min = 999, 100
	return fmt.Sprintf("unknown_file_%d.txt", rand.Intn(max-min)+min)
}

func createTmpBranchProtectionConfigFile() (file *os.File) {
	file, err := ioutil.TempFile("", "*.txt")
	if err != nil {
		log.Fatal(err)
	}
	return
}
