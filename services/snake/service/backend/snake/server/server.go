package server

import "net"

func Run() {
	addr := "0.0.0.0:7999"
	listener, err := net.Listen("tcp4", addr)
	if err != nil {
		return
	}
	defer listener.Close()
	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		gameConn := NewGameConn(conn)
		go gameConn.handleConnection()
	}
}
