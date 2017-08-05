package malScan

import (
	"github.com/hillu/go-yara"
	"io/ioutil"
	"log"
)

var rules = make([](*yara.Rules), 0)

func Scan(file []byte) []yara.MatchRule {
	for _, r := range rules {
		m, _ := r.ScanMem(file, 0, 0)
		if len(m) > 0 {
			return m
		}
	}
	return make([]yara.MatchRule, 0)
}

func ConstructRules() {
	dirName := "./rules/"
	fileList, _ := ioutil.ReadDir(dirName)
	for _, file := range fileList {
		filedata, _ := ioutil.ReadFile(dirName + file.Name())
		r, err := yara.Compile(string(filedata), nil)
		if err != nil {
			log.Printf("Error loading rules: %s", file.Name())
		} else {
			rules = append(rules, r)
		}
	}
}
