package util

import (
	"../proto"
	"fmt"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
)

var log = logrus.New()

func RequestLog(p proto.RequestParam) {
	log.WithFields(logrus.Fields{
		"Method": proto.MethodToStr(p.Method),
		"Size":   p.DataSize,
		"Hash":   fmt.Sprintf("%x", p.Hash),
	}).Info("Request: ")
}

func ResponseLog(p proto.ResponseParam) {
	log.WithFields(logrus.Fields{
		"Response Code": proto.ResponseToStr(p.Response),
		"Size":          p.DataSize,
	}).Info("Response: ")
}

func SaveLog(fileName string) {
	log.WithFields(logrus.Fields{
		"File Name": fileName,
	}).Info("Save File: ")
}

func LogFatal(message string, err error) {
	log.Fatal(message, err)
}

func init() {
	if os.Getenv("DEBUG") == "" {
		log.Out = ioutil.Discard
	} else {
		log.Out = os.Stdout
	}
	log.Formatter = new(logrus.TextFormatter)
}
