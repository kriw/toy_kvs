package main

import (
	"./server"
	"fmt"
	"log"
	"net/http"
	"net/http/pprof"
	"os"
)

func hiHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hi"))
}
func main() {
	go func() {
		r := http.NewServeMux()
		r.HandleFunc("/", hiHandler)

		// Register pprof handlers
		r.HandleFunc("/debug/pprof/", pprof.Index)
		r.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
		r.HandleFunc("/debug/pprof/profile", pprof.Profile)
		r.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
		r.HandleFunc("/debug/pprof/trace", pprof.Trace)
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()
	if len(os.Args) > 1 {
		fmt.Printf("Listening: %s\n", os.Args[1])
		server.Serve("tcp", os.Args[1])
	} else {
		fmt.Printf("Usage: %s <ip-address:port>\n", os.Args[0])
	}
}
