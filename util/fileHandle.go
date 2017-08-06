package util

import "strings"

func MatchExt(filename string, ext string) bool {
	ns := strings.Split(filename, ".")
	return ns[len(ns)-1] == ext
}
