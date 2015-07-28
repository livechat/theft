package main

import (
	// "io/ioutil"
	"log"
	"net/http"
	"time"
	"runtime"
	"html/template"
)

var logger *Logger

func main() {
	parseFlags()
	logger = &Logger{level:*settings.logLevel, path:settings.logPath}
	logger.Init()
	defer logger.Close()
	go goroutines()

	http.HandleFunc("/alive", alive)
	http.HandleFunc("/inspector/ws", handshake)
	http.HandleFunc("/hijacker/ws", handshake)
	http.HandleFunc("/hijacker/static", serveHijacker)

	fs := http.FileServer(http.Dir("frontend"))
	http.Handle("/inspector/", http.StripPrefix("/inspector/", fs))

	err := http.ListenAndServe(*settings.port, nil)
	if err != nil {
		log.Fatal("ERROR::", err)
	}
}

func serveHijacker(w http.ResponseWriter, r *http.Request){
	template, _ := template.ParseFiles("./hijack/hijack.js")
	w.Header().Set("Content-type", "application/javascript")
	url := ""

	if *settings.secure {
		url += "wss://"
	}else{
		url += "ws://"
	}

	url += *settings.domain
	url += *settings.port
	url += "/hijacker/ws"

	template.Execute(w, map[string] string {"url": url})
}

func goroutines() {
	logger.Debug("GOROUTINES", runtime.NumGoroutine())
	time.Sleep(5 * time.Second)
	goroutines()
}

func alive(w http.ResponseWriter, r *http.Request){
	w.WriteHeader(200)
	w.Write([]byte{})
}