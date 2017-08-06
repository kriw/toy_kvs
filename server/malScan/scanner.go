package malScan

import (
	"../../util"
	"fmt"
	"github.com/go-fsnotify/fsnotify"
	"github.com/hillu/go-yara"
	"io/ioutil"
	"log"
)

const (
	DIR = "./rules/"
	EXT = "yara"
)

var rules = make([](*yara.Rules), 0)
var watcher *fsnotify.Watcher

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
	fileList, _ := ioutil.ReadDir(DIR)
	for _, file := range fileList {
		fileName := file.Name()
		if !util.MatchExt(fileName, EXT) {
			continue
		}
		filedata, _ := ioutil.ReadFile(DIR + fileName)
		r, err := yara.Compile(string(filedata), nil)
		if err != nil {
			log.Printf("Error loading rules: %s", fileName)
		} else {
			rules = append(rules, r)
		}
	}
}

func RunRuleWatcher() {
	watcher, _ = fsnotify.NewWatcher()
	defer watcher.Close()
	watcher.Add(DIR)
	for {
		select {
		case event := <-watcher.Events:
			if util.MatchExt(event.Name, EXT) {
				handleWatchEvent(event.Op)
			}
		case err := <-watcher.Errors:
			fmt.Println("ERROR", err)
		}
	}
}

func handleWatchEvent(op fsnotify.Op) {
	switch op {
	case fsnotify.Write:
		ConstructRules()
	case fsnotify.Remove:
		ConstructRules()
	}
}
