package mNet

import "markV5/mFace"

func newMessage(routeID, connID string, data []byte) mFace.MMessage {
	m := &message{
		routeID: routeID,
		connID:  connID,
		data:    data,
	}
	return m
}

type message struct {
	routeID string
	connID  string
	data    []byte
}

func (m *message) RouteID() string {
	return m.routeID
}

func (m *message) ConnID() string {
	return m.connID
}

func (m *message) Marshal() []byte {
	return m.data
}

func (m *message) Unmarshal(data []byte) {
	m.data = data
}
