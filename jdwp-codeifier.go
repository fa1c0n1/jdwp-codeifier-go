package main

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"net"
)

type JDWPClient struct {
	Con net.Conn
	Id  int32
}

type PacketHeader struct {
	Length int32
	Id     int32
	Flags  byte
	Code   [2]byte
}

type JDWPPacket struct {
	Header PacketHeader
	Data   []byte
}

const COMMANDSET_VIRTUALMACHINE byte = 0x01
const COMMAND_VERSION byte = 0x01

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
	jdwp.Con = conn
	return true, nil
}

func (jdwp *JDWPClient) CloseConn() {
	if jdwp.Con != nil {
		defer jdwp.Con.Close()
	}
}

func (jdwp *JDWPClient) CreatePacket(commandSet byte, command byte, data []byte) *JDWPPacket {
	jdwpPacket := new(JDWPPacket)
	jdwp.Id += 1
	jdwpPacket.Header.Id = jdwp.Id
	jdwpPacket.Header.Flags = 0x00
	jdwpPacket.Header.Length = int32(len(data)) + 11
	jdwpPacket.Header.Code[0] = commandSet
	jdwpPacket.Header.Code[1] = command
	jdwpPacket.Data = data
	return jdwpPacket
}

func (jdwp *JDWPClient) GetVersionInfo() {
	reqPacket := jdwp.CreatePacket(COMMANDSET_VIRTUALMACHINE, COMMAND_VERSION, make([]byte, 0))
	buf := make([]byte, reqPacket.Header.Length)
	buffer := bytes.NewBuffer(buf)
	err := binary.Write(buffer, binary.BigEndian, &reqPacket.Header)
	if err != nil {
		panic(err)
	}
	fmt.Println("buf:", buf)

	sendData := buffer.Bytes()[buffer.Len()-len(buf):]

	_, err = jdwp.Con.Write(sendData)
	if err != nil {
		fmt.Println("GetVersionInfo: Failed to send data!")
		return
	}

	headBuf := make([]byte, 11)
	_, err = jdwp.Con.Read(headBuf)
	if err != nil {
		fmt.Println("GetVersionInfo: Failed to recv header!")
		return
	}
	fmt.Println("headerBuf:", headBuf)
	headerBuffer := bytes.NewBuffer(headBuf)
	replyPacketHeader := new(PacketHeader)
	err = binary.Read(headerBuffer, binary.BigEndian, replyPacketHeader)
	if err != nil {
		panic(err)
	}

	//fmt.Println(replyPacketHeader.Length)
	//fmt.Println(replyPacketHeader.Id)
	//fmt.Println(replyPacketHeader.Flags)
	//fmt.Println(replyPacketHeader.Code)

	replyDataBuf := make([]byte, replyPacketHeader.Length)
	_, err = jdwp.Con.Read(replyDataBuf)
	if err != nil {
		fmt.Println("GetVersionInfo: Failed to recv data!")
		return
	}

	fmt.Println(hex.Dump(replyDataBuf))
	fmt.Println(string(replyDataBuf))
}

func NewJDWPClient() *JDWPClient {
	return &JDWPClient{
		Id: 1,
	}
}

func main() {
	jdwpClient := NewJDWPClient()
	bRet, err := jdwpClient.connect("104.199.140.152", "8777")
	if !bRet {
		panic(err)
	}

	jdwpClient.GetVersionInfo()
}
