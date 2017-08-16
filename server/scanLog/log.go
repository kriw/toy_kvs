package scanLog

import (
	"../../proto"
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

var logger log.Logger
var logFile *os.File

func Write(ruleName string, hash [proto.HashSize]byte) {
	hashStr := fmt.Sprintf("%x", hash)
	logger.Printf("ruleName:%s, hash:%s\n", ruleName, hashStr)
}

func ChangeLogger() {
	logFile.Close()
	NewLogger()
}

func NewLogger() {
	fileName := strings.Replace(time.Now().String(), " ", "_", -1)
	fileName = strings.Replace(fileName, "-", "_", -1)
	fileName = fmt.Sprintf("logFiles/since_%s.log", fileName)
	logFile, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	logger.SetOutput(logFile)
}

func CloseLog() {
	logFile.Close()
}
