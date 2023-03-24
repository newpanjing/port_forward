package main

import (
	"flag"
	"fmt"
	"log"
	"net"
)

func main() {

	//解析参数
	var host = flag.String("listen", "0.0.0.0:9999", "listen address")
	var target = flag.String("target", "127.0.0.1:8899", "target address")
	flag.Parse()

	listen, err := net.Listen("tcp", *host)
	if err != nil {
		fmt.Println("Error listening:", err.Error())
	}
	defer func(listen net.Listener) {
		err := listen.Close()
		if err != nil {
			fmt.Println("Error closing listener:", err.Error())
		}
	}(listen)
	fmt.Println("Listening on ", *host, " forwarding to ", *target, " ...")

	for {
		conn, err := listen.Accept()
		if err != nil {
			fmt.Println("Error accepting:", err.Error())
			continue
		}
		go handleConnection(conn, *target)
	}
}

func handleConnection(conn net.Conn, targetAddr string) {
	defer func(conn net.Conn) {
		err := conn.Close()
		if err != nil {
			fmt.Println("Error closing connection:", err.Error())
		}
	}(conn)

	target, err := net.Dial("tcp", targetAddr)
	if err != nil {
		fmt.Println("Error dialing", targetAddr, err.Error())
		return
	}
	defer func(target net.Conn) {
		err := target.Close()
		if err != nil {
			fmt.Println("Error closing target connection:", err.Error())
		}
	}(target)
	fmt.Println("Connection established between", conn.RemoteAddr(), "and", targetAddr)
	go copyIO(conn, target)
	copyIO(target, conn)
}

func copyIO(src net.Conn, dst net.Conn) {
	buf := make([]byte, 1024)
	for {
		n, err := src.Read(buf)
		if err != nil {
			return
		}
		if n > 0 {
			log.Printf("read %d bytes from %s", n, src.RemoteAddr())
			_, err = dst.Write(buf[:n])
			if err != nil {
				return
			}
		}
	}
}
