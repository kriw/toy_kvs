package util

import (
	"../tkvsProtocol"
	log "github.com/sirupsen/logrus"
	"os"
)

var log = logrus.New()

func RequestLog(proto tkvsProtocol.RequestParam) {
	log.WithFields(log.Fields{
		"Method": tkvsProtocol.RequestToStr(proto.Method),
		"Size":   proto.DataSize,
		"Hash":   fmt.Sprintf("%x", proto.Hash),
	}).Info("A walrus appears")
}
func ResponseLog(proto tkvsProtocol.ResponseParam) {
	log.WithFields(log.Fields{
		"Response Code": tkvsProtocol.RequestToStr(proto.Response),
		"Size":          proto.DataSize,
		"Hash":          fmt.Sprintf("%x", proto.Hash),
	}).Info("A walrus appears")
}

func init() {
	log.Out = os.Stdout
	log.Formatter = new(logrus.TextFormatter)
}
