package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
)

const socks5VER = 0x05
const cmdBind = 0x01
const atypeIPV4 = 0x01
const atypeHOST = 0x03
const atypeIPV6 = 0x04

func main() {
	server, err := net.Listen("tcp", "127.0.0.1:1080")
	if err != nil {
		panic(err)
	}
	for {
		client, err := server.Accept()
		if err != nil {
			log.Printf("Accept failed %v", err)
			continue
		}
		go process(client)
	}
}

func auth(reader *bufio.Reader, conn net.Conn) (err error) {
	// +-----+---------+----------+
	// | VER | NMETHOD |  METHOD  |
	// +-----+---------+----------+
	// |  1  |    1    | 1 to 255 |
	// +-----|---------+----------+
	// VER:协议版本，socks5为 0x05
	// NMETHOD:支持认证的方法数量，为METHOD的字节数
	// METHOD:用来指示客户端和代理服务器之间的认证方法
	ver, err := reader.ReadByte()
	if err != nil {
		return fmt.Errorf("read ver failed:%w", err)
	}
	if ver != socks5VER {
		return fmt.Errorf("not supported ver:%v", ver)
	}
	methodSize, err := reader.ReadByte()
	if err != nil {
		return fmt.Errorf("read methodSize failed:%w", err)
	}
	method := make([]byte, methodSize)
	_, err = io.ReadFull(reader, method)
	if err != nil {
		return fmt.Errorf("read method failed:%w", err)
	}
	log.Println("ver", ver, "method", method)
	_, err = conn.Write([]byte{socks5VER, 0x00})
	if err != nil {
		return fmt.Errorf("write failed:%w", err)
	}
	return nil
}

func process(conn net.Conn) {
	defer conn.Close()
	reader := bufio.NewReader(conn)
	err := auth(reader, conn)
	if err != nil {
		log.Printf("Client %v auth failed:%v", conn.RemoteAddr(), err)
	}
	log.Println("Auth success")
}
