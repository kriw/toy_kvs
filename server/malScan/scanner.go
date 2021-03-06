package malScan

import (
	"../../util"
	"../scanLog"
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

var rules [](*yara.Rules)
var watcher *fsnotify.Watcher

func Scan(file []byte) []yara.MatchRule {
	matches := make([]yara.MatchRule, 0)
	for _, r := range rules {
		m, _ := r.ScanMem(file, 0, 0)
		if len(m) > 0 {
			matches = append(matches, m...)
		}
	}
	return matches
}

func ConstructRules() {
	rules = make([](*yara.Rules), 0)
	f := func(fileName string) {
		if util.MatchExt(fileName, EXT) {
			filedata, _ := ioutil.ReadFile(DIR + fileName)
			r, err := yara.Compile(string(filedata), nil)
			if err != nil {
				log.Printf("Error loading rules: %s", fileName)
			} else {
				rules = append(rules, r)
			}
		}
	}
	util.FilesMap(DIR, f)
}

func RunRuleWatcher() {
	watcher, _ = fsnotify.NewWatcher()
	defer watcher.Close()
	scanLog.NewLogger()
	defer scanLog.CloseLog()
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
		fallthrough
	case fsnotify.Remove:
		ConstructRules()
		scanLog.ChangeLogger()
	}
}
