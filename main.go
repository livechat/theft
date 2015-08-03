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
	staticFileServerHandler := http.StripPrefix("/inspector/", fs);
	auth := Base64Authorization{&staticFileServerHandler}

	http.Handle("/inspector/", auth)

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

type Base64Authorization struct {
	handler *http.Handler
}

func (self Base64Authorization) ServeHTTP(w http.ResponseWriter, r *http.Request){
	if *settings.auth != "" {
		username, password, ok := r.BasicAuth()

		if *settings.auth != username + ":" + password || ! ok {
			w.Header().Set("WWW-Authenticate", "Basic")
			w.WriteHeader(401)
			w.Write([]byte{})
			return
		}
	}

	(*self.handler).ServeHTTP(w, r)
}