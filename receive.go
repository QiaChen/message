package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

var websocketService *hub

func receive() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
			return
		}
	}()
	startHttp()
}

func startHttp() {
	websocketService = newHub()
	go websocketService.run()
	http.HandleFunc("/send", sendHandler)           //设置访问的路由
	http.HandleFunc("/subscribe", subscribeHandler) //设置访问的路由
	http.Handle("/ws", wsHandler{h: websocketService})
	err := http.ListenAndServe(":9091", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func subscribeHandler(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Fprintf(w, "{\"code\":\"2\",\"msg\":\"系统错误\"}")
			return
		}
	}()
	r.ParseForm()
	sign := r.Form["sign"][0]
	uid := r.Form["uid"][0]
	getTopics := r.Form["topics"][0]

	if smd5(getTopics+"ajsdlkfjdslk@sd~!"+uid) != sign {
		fmt.Fprintf(w, "{\"code\":\"1\",\"msg\":\"sign verify fail\"}")
		return
	}
	arr := strings.Split(getTopics, ",")
	for _, t := range arr {
		ListenTopic(uid, t)
	}
	fmt.Fprintf(w, "{\"code\":\"0\",\"Topics\":\""+getTopics+"\"}")
}
func sendHandler(w http.ResponseWriter, r *http.Request) {
	// fmt.Println("send-start")
	defer func() {
		if err := recover(); err != nil {
			//fmt.Println(r.Form)
			fmt.Fprintf(w, "{\"code\":\"2\",\"msg\":\"系统错误\"}")
			return
		}
	}()
	r.ParseForm()
	sign := r.Form["sign"][0]
	timestr := r.Form["time"][0]

	if smd5(timestr+"ajsdlkfjdslk@sd~!") != sign {
		fmt.Fprintf(w, "{\"code\":\"1\",\"msg\":\"sign verify fail\"}")
		return
	}

	str := r.Form["data"][0]
	// fmt.Println(time.Now().Format("2006-01-02 15:04:05")+"??"+str)
	send(str)
	fmt.Fprintf(w, "{\"code\":\"0\"}")
}

func newHub() *hub {
	return &hub{
		Broadcast:   make(chan []byte),
		register:    make(chan *connection),
		unregister:  make(chan *connection),
		connections: make(map[string]map[string]*connection),
	}
}

func (h *hub) run() {
	for {
		select {
		case c := <-h.register:
			// fmt.Println("connected")
			if _, ok := h.connections[c.uid]; ok {
				c.cid = strconv.Itoa(int(time.Now().Unix())) + strconv.Itoa(len(h.connections[c.uid]))
				h.connections[c.uid][c.cid] = c
			} else {
				//TODO 用户上线通知
				h.connections[c.uid] = make(map[string]*connection)
				c.cid = strconv.Itoa(int(time.Now().Unix())) + "0"
				h.connections[c.uid][c.cid] = c
			}
			//订阅私人频道
			ListenTopic(c.uid, "ToUser://"+c.uid)
		case c := <-h.unregister:
			if _, ok := h.connections[c.uid]; ok {
				// fmt.Println("stop")
				if numconnect := len(h.connections[c.uid]); numconnect > 1 {
					delete(h.connections[c.uid], c.cid)
				} else {
					//TODO 用户下线通知
					UnsubscribeUserAll(c.uid)
					delete(h.connections[c.uid], c.cid)
					delete(h.connections, c.uid)

				}
				close(c.send)
				c.ws.Close()
			}

		case m := <-h.Broadcast:
			fmt.Println("msg:" + string(m))
		}
	}
}

func (c *connection) reader() {
	for {
		_, message, err := c.ws.ReadMessage()
		if err != nil {
			break
		}

		msg := string(message)
		if c.uid == "cmopApi" {
			send(msg)
		} else {
			func(msg string) {}(msg)
		}

	}
	c.ws.Close()
}

func (c *connection) writer() {
	for message := range c.send {
		err := c.ws.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			break
		}
	}
	c.ws.Close()
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  4096,
	WriteBufferSize: 4096,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (wsh wsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	uid := r.Form.Get("uid")
	topics := r.Form.Get("topics")
	arr := strings.Split(topics, ",")
	for _, t := range arr {
		ListenTopic(uid, t)
	}

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	c := &connection{send: make(chan []byte, 256), ws: ws, h: wsh.h, uid: uid}
	c.h.register <- c
	defer func() { c.h.unregister <- c }()
	go c.writer()
	c.reader()
}
