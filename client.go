package main

import (
	"math/rand"
)

type Endpoint interface {
	GetId() int64
	GetKind() int64
	Unregister()
	Send ([]byte)
}

type Client struct {
	id int64
	conn *Connection
	kind int64
}

func (self *Client) GetId() (int64){
	return self.id
}

func (self *Client) GetKind() (int64) {
	return self.kind
}

func (self *Client) Unregister() {
	hub.unregister(self.id)
}

func (self *Client) Send(raw []byte) {
	self.conn.send(raw)
}

func NewClient(conn *Connection, kind int64) (*Client) {
	id := rand.Int63() / 1000
	client := Client{id, conn, kind} 

	return &client
}