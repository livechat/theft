package main

type Inspector struct {
	*Client

	inspected int64
}

func (self *Inspector) handle() {
	self.inspected = -1
	self.sendHijackersList()

	for {
		select {
			case raw, ok := <- self.conn.rx:
				if ok == false {
					self.removeListener()
					self.Unregister()
					return 
				}

				self.protocol(raw)
		}
	}
}

func (self *Inspector) sendHijackersList() {
	hijackers := hub.getEndpointsByKind(HIJACKER)
	JsonHijackers := []*JsonHijacker{}
	frame := Frame{ Event: "hijackers" }

	for _, endpoint := range hijackers {
		hijacker := (*endpoint).(*Hijacker)
		JsonHijackers = append(JsonHijackers, hijacker.getJsonHijacker())
	}

	frame.SetData(JsonHijackersEvent{JsonHijackers})
	hub.send(self.id, self.id, frame.GetRaw())
}

func (self *Inspector) removeListener() {
	if self.inspected != -1 {
		endpoint := hub.getEndpointsById(self.inspected);

		if (*endpoint) != nil {
			(*endpoint).(*Hijacker).unregisterListener(self.id)
		}

		self.inspected = -1
	}
}

func (self *Inspector) protocol (raw []byte){
	var (
		frame *Frame
		err error
	)

	if frame, err = CrateFrameFromRaw(raw); err != nil {
		return
	}

	switch frame.Event {
		case "inspect":
			inspect := JsonInspectEvent{}
			frame.GetData(&inspect)

			self.removeListener()

			endpoint := hub.getEndpointsById(inspect.Session);
			if endpoint == nil {
				return
			}

			self.inspected = (*endpoint).GetId()
			(*endpoint).(*Hijacker).registerListener(self.id)

		case "command":
			command := JsonCommand{}
			frame.GetData(&command)
			command.InspectorId = self.id

			frame.SetData(command)
			if command.Batch {
				hub.broadcast(self.id, frame.GetRaw())
			}else if (self.inspected != -1){
				hub.send(self.id, self.inspected, frame.GetRaw())
			}

			

		default:
			logger.Error("INSPECTOR", "::PROTOCOL", "missing command", frame.Event)
	}
}