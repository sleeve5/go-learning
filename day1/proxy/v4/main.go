package main

import (
	"bufio"
	"context"
	"encoding/binary"
	"errors"
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

func process(conn net.Conn) {
	defer conn.Close()
	reader := bufio.NewReader(conn)
	err := auth(reader, conn)
	if err != nil {
		log.Printf("Client %v auth failed:%v", conn.RemoteAddr(), err)
	}
	err = connect(reader, conn)
	if err != nil {
		log.Printf("client %v auth failed:%v", conn.RemoteAddr(), err)
		return
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
	// +----+--------+
	// |VER | METHOD |
	// +----+--------+
	// | 1  |   1    |
	// +----+--------+

	_, err = conn.Write([]byte{socks5VER, 0x00})
	if err != nil {
		return fmt.Errorf("write failed:%w", err)
	}
	return nil
}

func connect(reader *bufio.Reader, conn net.Conn) (err error) {
	// +----+-----+-------+-------+----------+----------+
	// |VER | CMD |  RSV  | ATYPE | DST.ADDR | DST.PORT |
	// +----+-----+-------+-------+----------+----------+
	// | 1  |  1  | X'00' |   1   | Variable |    2     |
	// +----+-----+-------+-------+----------+----------+
	// VER:版本号，socks5的值为0x05
	// CMD:0x01表示CONNECT请求
	// RSV:保留字段，值为0x00
	// ATYPE:目标地址类型，DST.ADDR的数据对应这个字段的类型。
	// 		0x01表示IPv4地址，DST.ADDR为4个字节
	// 		0x03表示域名，DST.ADDR是一个可变长度的域名
	// DST.ADDR:一个可变长度的值
	// DST.PORT:目标端口，固定2个字节

	buf := make([]byte, 4)
	_, err = io.ReadFull(reader, buf)
	if err != nil {
		return fmt.Errorf("read header failed:%w", err)
	}
	ver, cmd, atype := buf[0], buf[1], buf[3]
	if ver != socks5VER {
		return fmt.Errorf("not supported version:%w", err)
	}
	if cmd != cmdBind {
		return fmt.Errorf("not supported cmd:%w", err)
	}
	addr := ""
	switch atype {
	case atypeIPV4:
		_, err = io.ReadFull(reader, buf)
		if err != nil {
			return fmt.Errorf("read atype failed:%w", err)
		}
		addr = fmt.Sprintf("%d.%d.%d.%d", buf[0], buf[1], buf[2], buf[3])
	case atypeHOST:
		hostSize, err := reader.ReadByte()
		if err != nil {
			return fmt.Errorf("read hostsize failed:%w", err)
		}
		host := make([]byte, hostSize)
		_, err = io.ReadFull(reader, host)
		if err != nil {
			return fmt.Errorf("read host failed:%w", err)
		}
		addr = string(host)
	case atypeIPV6:
		return errors.New("IPV6 is not supported")
	default:
		return errors.New("invalid atype")
	}
	_, err = io.ReadFull(reader, buf[:2])
	if err != nil {
		return fmt.Errorf("read port failed:%w", err)
	}
	port := binary.BigEndian.Uint16(buf[:2])

	dest, err := net.Dial("tcp", fmt.Sprintf("%v:%v:", addr, port))
	if err != nil {
		return fmt.Errorf("dial dest failed:%w", err)
	}
	defer dest.Close()
	log.Println("dial", addr, port)

	// +----+-----+-------+-------+----------+----------+
	// |VER | REP |  RSV  | ATYPE | BND.ADDR | BND.PORT |
	// +----+-----+-------+-------+----------+----------+
	// | 1  |  1  | X'00' |   1   | Variable |    2     |
	// +----+-----+-------+-------+----------+----------+
	// VER socks版本，这里为0x05
	// REP Relay field,内容取值如下 X’00’ succeeded
	// RSV 保留字段
	// ATYPE 地址类型
	// BND.ADDR 服务绑定的地址
	// BND.PORT 服务绑定的端口DST.PORT

	_, err = conn.Write([]byte{0x05, 0x00, 0x00, 0x01, 0, 0, 0, 0, 0, 0})
	if err != nil {
		return fmt.Errorf("write failed:%w", err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		_, _ = io.Copy(dest, reader)
		cancel()
	}()
	go func() {
		_, _ = io.Copy(conn, dest)
		cancel()
	}()

	<-ctx.Done()
	return nil
}
