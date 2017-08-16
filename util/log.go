package util

import (
	"../proto"
	"fmt"
	"github.com/sirupsen/logrus"
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

func logFatal(err error) {
	log.Fatal(err)
}

func init() {
	log.Out = os.Stdout
	log.Formatter = new(logrus.TextFormatter)
}
