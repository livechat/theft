package main

import (
  "sync"
)

const (
		HIJACKER = 1 << iota
		INSPECTOR = 1 << iota
)

type Hub struct {
	*sync.RWMutex
	endpoints map[int64]Endpoint
}

func (self *Hub) register(conn *Connection) {

	switch conn.path {
		case "/hijacker/ws":
			hijacker := Hijacker{Client:NewClient(conn, HIJACKER)}
			go hijacker.handle()

		default:
			inspector := Inspector{Client:NewClient(conn, INSPECTOR)}
			self.addEndpoint(inspector)
			go inspector.handle()
	}
}

func (self *Hub) unregister(id int64){
	self.Lock()
	defer self.Unlock()
	delete(self.endpoints, id)
}

func (self *Hub) addEndpoint(endpoint Endpoint) {
	self.Lock()
	defer self.Unlock()

	self.endpoints[endpoint.GetId()] = endpoint
}

func (self *Hub) getEndpointsById(id int64) *Endpoint {
	self.RLock()
	defer self.RUnlock()

	endpoint := self.endpoints[id]
	return &endpoint
}

func (self *Hub) getEndpointsByKind(kind int64) []*Endpoint {
	self.RLock()
	defer self.RUnlock()

	var endpoints []*Endpoint

	for key, _ := range self.endpoints {
		endpoint := self.endpoints[key]

		if endpoint.GetKind() == kind {
			endpoints = append(endpoints, &endpoint)
		}
	}

	return endpoints
}

func (self *Hub) broadcast (from int64, raw []byte){
	self.RLock()
	defer self.RUnlock()
	endpoint := self.endpoints[from]
	if endpoint == nil {
		return 
	}

	kind := endpoint.GetKind()

	for _, endpoint := range self.endpoints {
		if endpoint.GetKind() != kind {
			endpoint.Send(raw)
		}
	}
}

func (self *Hub) send (from int64, to int64, raw []byte) {
	self.RLock()
	defer self.RUnlock()

	endpoint := self.endpoints[to]
	if endpoint == nil {
		return
	}

	self.endpoints[to].Send(raw)
}

var hub = Hub{&sync.RWMutex{}, make(map[int64]Endpoint)}