package main

import (
	"errors"
	"fmt"
	"net"
)

type JDWPClient struct {
}

type PacketHeader struct {
	Length int
	Id     int
	Flags  byte
	Code   [2]byte
}

type JDWPPacket struct {
	Header PacketHeader
	Data   []byte
}

const (
	JDWP_HANDSHAKE = "JDWP-Handshake"
)

func (jdwp *JDWPClient) connect(host string, port string) (bool, error) {
	targetAddr := net.JoinHostPort(host, port)
	fmt.Println("targetAddr:", targetAddr)

	conn, err := net.Dial("tcp", targetAddr)
	if err != nil {
		fmt.Println(err)
		return false, err
	}

	defer conn.Close()

	n, err := conn.Write([]byte(JDWP_HANDSHAKE))
	if err != nil {
		fmt.Println(err)
		return false, err
	}

	fmt.Println("send, n=", n)
	readBs := make([]byte, len(JDWP_HANDSHAKE))
	n, err = conn.Read(readBs)
	if n != len(JDWP_HANDSHAKE) {
		return false, errors.New("Fail to jdwp handshake!!!")
	}

	fmt.Println("recv:", string(readBs[:n]))
	fmt.Println("JDWP Handshake successfully.")
	return true, nil
}

func (jdwp *JDWPClient) CreatePacket() {
	fmt.Println("Create Packet")
}

func main() {
	jdwpClient := new(JDWPClient)
	bRet, err := jdwpClient.connect("104.199.140.152", "8777")
	if !bRet {
		panic(err)
	}
}
