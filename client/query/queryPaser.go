package query

import (
	"strings"
)

type QueryMethod int

const (
	GET QueryMethod = iota
	SET
	SETFILE
	Unknown
)

type Query struct {
	Op   QueryMethod
	Args [][]byte
}

func parseOp(op string) QueryMethod {
	switch op {
	case "get":
		return GET
	case "set":
		return SET
	case "setfile":
		return SETFILE
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
	op, argsStr := strs[0], strs[1:]
	args := make([][]byte, 0)
	for _, a := range argsStr {
		args = append(args, []byte(a))
	}
	return Query{parseOp(op), args}
}
