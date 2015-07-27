package main

import (
	"time"
	"sync"
)

type Delay struct {
	delay int64
	start int64
	mutex *sync.Mutex
}

func (self *Delay) Mark() {
	self.mutex.Lock()
	defer self.mutex.Unlock()

	self.start = time.Now().UnixNano();
}

func (self *Delay) Ready() {
	self.mutex.Lock()
	defer self.mutex.Unlock()
	self.delay = time.Now().UnixNano() - self.start;
}

func (self *Delay) GetMicroDelay() int64 {
	self.mutex.Lock()
	defer self.mutex.Unlock()

	return self.delay/1000;
}