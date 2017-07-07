package query

import (
	"strings"
)

type QueryKind int

const (
	GET QueryKind = iota
	SET
	Unknown
)

type Query struct {
	Op   QueryKind
	Args []string
}

func parseOp(op string) QueryKind {
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
