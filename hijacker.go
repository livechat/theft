package main

import (
	"time"
	"sync"
)

type Hijacker struct {
	*Client
	*sync.RWMutex

	userAgent string
	location string

	listeners map[int64]bool
}

func (self *Hijacker) init() {
	self.RWMutex = &sync.RWMutex{}
	self.listeners = make(map[int64]bool)
}

func (self *Hijacker) handle() {
	self.init()

	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for {
		select {
			case frame, ok := <- self.conn.rx:
				if ok == false {
					self.sendHijackerEventFrame("unregister")
					self.Unregister()
					return
				}

				self.protocol(frame)

			case <-ticker.C:
				self.sendHijackerEventFrame("delay")
		}
	}
} 

func (self *Hijacker) registerListener(listener int64){
	self.Lock()
	defer self.Unlock()

	self.listeners[listener] = true;
	
}

func (self *Hijacker) unregisterListeners(){
	self.Lock();
	defer self.Unlock()

	for listener, _ := range(self.listeners){
		delete(self.listeners, listener)
	}
}

func (self *Hijacker) unregisterListener(listener int64){
	self.Lock()
	defer self.Unlock()
	delete(self.listeners, listener)
}

func (self *Hijacker) sendHijackerEventFrame(kind string) {
	jsonHijackerEvent := JsonHijackerEvent{kind, self.getJsonHijacker()}
	frame := Frame{Event:"hijacker"}
	frame.SetData(jsonHijackerEvent)
	hub.broadcast(self.id, frame.GetRaw());
}

func (self *Hijacker) getJsonHijacker() *JsonHijacker {
	jsonHijacker := &JsonHijacker{self.id, self.userAgent, self.location, self.conn.delay.GetMicroDelay()}
	return jsonHijacker
}

func (self *Hijacker) protocol (raw []byte){
	var (
		frame *Frame
		err error
	)

	if frame, err = CrateFrameFromRaw(raw); err != nil {
		return
	}

	switch frame.Event {
		case "info":
			info := JsonHijacker{}
			frame.GetData(&info)
			self.id = info.Session
			self.userAgent = info.Browser
			self.location = info.Location

			hub.addEndpoint(self)
			self.sendHijackerEventFrame("register")

		case "log":
			self.RLock()
			defer self.RUnlock()

			for listener, _ := range(self.listeners){
				hub.send(self.id, listener, frame.GetRaw())
			}
			
		default:
			logger.Error("HIJACKER", "::PROTOCOL", "missing command", frame.Event)

	}
}