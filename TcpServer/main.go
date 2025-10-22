package main

import (
	"log"
	"net"
	"time"
)

func handleConnection(conn net.Conn) {
	log.Println(conn.RemoteAddr())
	var buf []byte = make([]byte, 1024)
	_, err := conn.Read(buf)
	if err != nil {
		log.Fatal(err)
	}
	time.Sleep(time.Second * 10)
	conn.Write([]byte("HTTP/1.1 200 OK\r\n\r\nHello, world\r\n"))
	conn.Close()
}
func main() {
	listener, err := net.Listen("tcp", ":6379")
	if err != nil {
		log.Fatal(err)
	}
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal(err)
		}
		go handleConnection(conn)
	}
}
