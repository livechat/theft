package main

import (
	"encoding/json"
)

type Frame struct {
	Event string `json:"event"`
	Data json.RawMessage `json:"data"`
}

type JsonLogEvent struct {
	Session int64 `json:"session"`
	Log string `json:"log"`
}

type JsonHijackersEvent struct {
	Hijackers []*JsonHijacker `json:"hijackers"`
}

type JsonHijackerEvent struct {
	Kind string `json:"kind"`
	Hijacker *JsonHijacker `json:"hijacker"`
}

type JsonHijacker struct {
	Session int64 `json:"session"`
	Browser string `json:"browser"`
	Location string `json:"location"`
	Delay int64 `json:"delay"`
}

type JsonDelayEvent struct {
	Session int64 `json:"session"`
	Delay int64 `json:"delay"`
}

type JsonInspectEvent struct {
	Session int64 `json:"session"`
}

type JsonCommand struct {
	Id int64 `json:"id"`
	HijackerId int64 `json:"hijacker_id"`
	InspectorId int64 `json:"inspector_id"`
	Cmd string `json:"cmd"`
	Response string `json:"response"`
	Batch bool `json:"batch"`
	Echo bool `json:"echo"`
}

func (self *Frame) GetRaw() []byte {
	raw, _ := json.Marshal(self)
	return raw
}

func (self *Frame) SetData(value interface{}) {
	raw, _ := json.Marshal(value)
	self.Data = raw
}

func (self *Frame) GetData(value interface{}) {
	json.Unmarshal(self.Data, value)
}

func CrateFrameFromRaw(raw []byte) (*Frame, error) {
	frame := &Frame{}

	if err := json.Unmarshal(raw, frame); err != nil {
		return nil, err
	}

	return frame, nil
}