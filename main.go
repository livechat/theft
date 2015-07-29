package main

import (
	"log"
	"net/http"
	"time"
	"runtime"
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
	http.HandleFunc("/hijacker/static", serveHijackerClient)

	fs := http.FileServer(http.Dir("frontend"))
	http.Handle("/inspector/", http.StripPrefix("/inspector/", fs))

	err := http.ListenAndServe(*settings.port, nil)
	if err != nil {
		log.Fatal("ERROR::", err)
	}
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