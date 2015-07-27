package main

import (
	"net/http"
	"github.com/gorilla/websocket"
	"time"
	"math/rand"
	"sync"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool { return true },
}

type Connection struct {
	*sync.Mutex

	id int64
	ws *websocket.Conn
	tx chan []byte
	rx chan []byte

	closed bool

	delay Delay
	path string
}

func (self *Connection) close() {
	self.Lock();
	defer self.Unlock();

	if self.closed == false {
		self.ws.Close()

		close(self.tx)
		close(self.rx)

		self.closed = true;
	}	
}

func (self *Connection) read() {
	defer func(){
		self.close()
		logger.Debug("CONNECTION", "::RX", "CLOSED");
	}()

	self.ws.SetReadDeadline(time.Now().Add(60 * time.Second))
	self.ws.SetPongHandler(func(string) error { 
		self.delay.Ready()
		self.ws.SetReadDeadline(time.Now().Add(60 * time.Second)); 
		logger.Debug("CONNECTION", "::PONG", self.id, self.delay.GetMicroDelay());
		return nil 
	})

	for {
		_, message, err := self.ws.ReadMessage()
		if err != nil {
			logger.Debug("CONNECTION", "::RX", "::ERR", self.id, err);
			return
		}

		logger.Debug("CONNECTION", "::RX", self.id, string(message[:]));
		self.rx <- message
	}
}

func (self *Connection) write() {
	ticker := time.NewTicker(10 * time.Second)
	defer func() {
		ticker.Stop()
		self.close()
		logger.Debug("CONNECTION", "::TX", "CLOSED");
	}()

	for {
		select {
			case <-ticker.C:
				self.ws.SetWriteDeadline(time.Now().Add(10 * time.Second))
				self.delay.Mark()
				if err := self.ws.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
					logger.Debug("CONNECTION", "::PING", self.id, err);
					return
				}
				logger.Debug("CONNECTION", "::PING", self.id);

			case frame := <- self.tx:
				self.ws.SetWriteDeadline(time.Now().Add(10 * time.Second))
				if err := self.ws.WriteMessage(websocket.TextMessage, frame); err != nil {
					return
				}
				logger.Debug("CONNECTION", "::TX", self.id, string(frame[:]));
		}
	}
}

func (self *Connection) send(frame []byte) {
	self.Lock()
	defer self.Unlock()

	if self.closed == false {
		self.tx <- frame
	}		
}

func handshake(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", 405)
		return
	}

	ws, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		logger.Warning("CONNECTION", "::UPGRAGE", err);
		return
	}

	self := &Connection{
		&sync.Mutex{},
		rand.Int63n(900000) + 100000, 
		ws, 
		make(chan []byte), 
		make(chan []byte), 
		false, 
		Delay{delay:0, mutex: &sync.Mutex{}},
		r.URL.Path }

	logger.Debug("CONNECTION", "::NEW", self.id);
	go self.read()
	go self.write()

	hub.register(self)
}
