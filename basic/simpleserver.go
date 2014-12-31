/*
Skills: simple webserver, simple fileserver with stripprefix, JSON response
*/
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"runtime"
)

func main() {
	portPtr := flag.String("port", "10000", "port for listening")
	pathPtr := flag.String("path", ".", "Path for fileserver")
	flag.Parse()

	port := ":" + *portPtr
	fmt.Printf("Go to http://127.0.0.1%s/stopall for Stop and exit.\n", port)
	fmt.Printf("Go to http://127.0.0.1%s/sinfo for Stop and exit.\n", port)
	http.Handle("/", http.StripPrefix("/", http.FileServer(http.Dir(*pathPtr))))

	http.HandleFunc("/stopall", func(w http.ResponseWriter, r *http.Request) {
		log.Fatal("Stopped by user")
	})
	http.HandleFunc("/sinfo", func(w http.ResponseWriter, r *http.Request) {
		type ServerInfo struct {
			Version       string `json:"version"`
			NumGoroutines int    `json:"goroutines"`
			NumCPU        int    `json:"cpu"`
		}
		group := ServerInfo{
			Version:       runtime.Version(),
			NumGoroutines: runtime.NumGoroutine(),
			NumCPU:        runtime.NumCPU(),
		}
		b, err := json.Marshal(group)
		if err != nil {
			w.Write([]byte(err.Error()))
		} else {
			w.Write(b)
		}
	})

	log.Println("Simple server running!")
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
