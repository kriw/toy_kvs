package query

import (
	"strings"
)

type QueryMethod int

const (
	GET QueryMethod = iota
	SET
	Unknown
)

type Query struct {
	Op   QueryMethod
	Args []string
}

func parseOp(op string) QueryMethod {
	switch op {
	case "get":
		return GET
	case "set":
		return SET
	default:
		return Unknown
	}
}

func trimEach(strs []string) []string {
	for i, str := range strs {
		strs[i] = strings.TrimSpace(str)
	}
	return strs
}

func Parse(queryStr string) Query {
	strs := strings.Split(queryStr, " ")
	strs = trimEach(strs)
	op, args := strs[0], strs[1:]
	return Query{parseOp(op), args}
}
