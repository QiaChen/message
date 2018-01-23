package main

import (
	"bufio"
	"fmt"
	"net"
	"time"
)

const (
	Head = 4
)

var (
	ClientMap map[int]net.Conn = make(map[int]net.Conn)
)

func CrossDomainServer() {

	tcpAddr, err := net.ResolveTCPAddr("tcp4", ":843")
	checkError(err)
	listener, err := net.ListenTCP("tcp", tcpAddr)
	checkError(err)
	clientIndex := 0

	for {
		clientIndex++

		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		go handleClient(conn, clientIndex)
	}
}

func handleClient(conn net.Conn, index int) {
	ClientMap[index] = conn
	fc := func() {
		time.Sleep(time.Second)
		conn.Close()
		delete(ClientMap, index)

	}
	defer fc()
	sendFirstMsg(conn)
}
func sendFirstMsg(conn net.Conn) {
	str := `<?xml version="1.0"?>  
            <!DOCTYPE cross-domain-policy SYSTEM "/xml/dtds/cross-domain-policy.dtd">  
            <cross-domain-policy>  
                <site-control permitted-cross-domain-policies="master-only"/>  
                <allow-access-from domain="*" to-ports="*" />  
            </cross-domain-policy>`
	writer := bufio.NewWriter(conn)
	writer.WriteString(str)
	writer.Flush()
}

func checkError(err error) {
	if err != nil {
		fmt.Println(err)
	}
}
