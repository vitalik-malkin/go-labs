package main

import (
	"bytes"
	"encoding/base64"
	"flag"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

const (
	base64Alphabet = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/"

	base64ActionName = "BASE64"
)

var (
	actionNameFlag *string

	inputPathFlag  *string
	outputPathFlag *string
)

func main() {
	actionNameFlag = flag.String("a", "", "action to execute")
	inputPathFlag = flag.String("i", "", "input path")
	outputPathFlag = flag.String("o", "", "output path")
	flag.Parse()

	actionName := *actionNameFlag

	var actionErr error

	switch strings.ToUpper(actionName) {
	case base64ActionName:
		r, err := contentAsBase64(*inputPathFlag)
		if err != nil {
			actionErr = err
			break
		}
		err = ioutil.WriteFile(*outputPathFlag, []byte(r), 0644)
		if err != nil {
			actionErr = err
		}
	default:
		log.Fatalf("unknown action specified: '%v'", actionName)
	}
	if actionErr != nil {
		log.Fatalf("error while executing action '%s':\n%v", actionName, actionErr)
	}
}

func contentAsBase64(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	fileContent, err := ioutil.ReadAll(file)
	if err != nil {
		return "", err
	}
	buf := new(bytes.Buffer)
	encoding := base64.NewEncoding(base64Alphabet)
	encoder := base64.NewEncoder(encoding, buf)
	_, err = encoder.Write(fileContent)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}
