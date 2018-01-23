package main

import (
	"github.com/gorilla/websocket"
)

type hub struct {
	connections map[string]map[string]*connection
	Broadcast   chan []byte
	register    chan *connection
	unregister  chan *connection
}

type wsHandler struct {
	h *hub
}

type connection struct {
	ws   *websocket.Conn
	uid  string
	cid  string
	send chan []byte
	h    *hub
}

type requestMap struct {
	SendType string
	MsgType  string
	Msg      string
	Topic    string
	Uid      string
	Fuid     string
	Data     map[string]string
}
type topic struct {
	Name  string
	Users map[string]string
}
