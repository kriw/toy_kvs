package util

import (
	"io/ioutil"
	"strings"
)

func MatchExt(filename string, ext string) bool {
	ns := strings.Split(filename, ".")
	return ns[len(ns)-1] == ext
}

func FilesMap(dirName string, f func(string)) {
	fileList, _ := ioutil.ReadDir(dirName)
	for _, file := range fileList {
		fileName := file.Name()
		f(fileName)
	}
}
