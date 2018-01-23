package main

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"strings"
)

func send(str string) {

	reqMap := requestMap{}
	err := json.Unmarshal([]byte(str), &reqMap)
	if err != nil {
		return 
	}
	if reqMap.Topic == "" {
		reqMap.Topic = "ToUser://" + reqMap.Uid
	}
	types := strings.Split(reqMap.SendType, ",")
	for _, value := range types {
		if value == "websocket" {
			sendWebsocket(reqMap)
		}
	}
}

func sendWebsocket(notimap requestMap) {

	thisTop, err := Topics[strings.ToUpper(notimap.Topic)]
	if !err {
		return
	}
	for _, v := range thisTop.Users {
		users, err := websocketService.connections[v]
		if !err {
			continue
		}
		str, _ := json.Marshal(notimap)
		for _, c := range users {
			c.ws.WriteMessage(websocket.TextMessage, str)
		}
	}

}
