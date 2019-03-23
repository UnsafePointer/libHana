// package name: libHana
package main

/*
#cgo CFLAGS:
typedef unsigned int* (*GlobalRegistersCallbackType)();
unsigned int* callGlobalRegistersCallback(GlobalRegistersCallbackType callback);
*/
import "C"

import (
	"bufio"
	"bytes"
	"fmt"
	"net"
	"strings"
	"unsafe"
)

type CallbackType func(C.uint)

type ScannerState byte

const (
	Beginning ScannerState = '$'
	End       ScannerState = '#'
)

var (
	state = Beginning
)

var globalRegistersCallback C.GlobalRegistersCallbackType
var ackDisabled = false

//export SetGlobalRegistersCallback
func SetGlobalRegistersCallback(fn C.GlobalRegistersCallbackType) {
	globalRegistersCallback = fn
}

//export StartDebugServer
func StartDebugServer(port uint) {
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", port))
	if err != nil {
		fmt.Println("Error listening: ", err.Error())
	}
	defer lis.Close()
	conn, err := lis.Accept()
	if err != nil {
		fmt.Println("Error accepting: ", err.Error())
	}
	for {
		handleConnection(conn)
	}
}

func scan(data []byte, atEOF bool) (int, []byte, error) {
	for i := 0; i < len(data); i++ {
		if data[i] == byte(state) {
			if state == Beginning {
				state = End
			} else {
				state = Beginning
			}
			return i + 1, data[:i], nil
		}
	}
	return 0, data, bufio.ErrFinalToken
}

func validatePacket(packet string, checksum string) bool {
	return true
}

func ack(conn net.Conn) {
	if ackDisabled {
		return
	}
	conn.Write([]byte("+"))
}

func nack(conn net.Conn) {
	conn.Write([]byte("-"))
}

func handleConnection(conn net.Conn) {
	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	if err != nil {
		return
	}
	scanner := bufio.NewScanner(bytes.NewReader(buf[:n]))
	scanner.Split(scan)
	for scanner.Scan() {
		packet := scanner.Text()
		switch packet {
		case "":
		case "+":
		case "-":
		default:
			scanner.Scan()
			checkSum := scanner.Text()
			if validatePacket(packet, checkSum) {
				ack(conn)
				reply(conn, packet)
			} else {
				nack(conn)
			}
		}
	}
}

func send(conn net.Conn, message string) {
	checksum := uint8(0)
	for i := 0; i < len(message); i++ {
		char := message[i]
		checksum += char
	}
	reply := fmt.Sprintf("$%s#%02x", message, checksum)
	conn.Write([]byte(reply))
}

func byteToString(b uint8) string {
	hexValues := [16]rune{'0', '1', '2', '3', '4', '5', '6', '7', '8', '9', 'a', 'b', 'c', 'd', 'e', 'f'}
	var msg strings.Builder
	msg.WriteRune(hexValues[b>>4])
	msg.WriteRune(hexValues[b&0xF])
	return msg.String()
}

func wordToString(word uint32) string {
	var msg strings.Builder
	for i := uint8(0); i < 4; i++ {
		b := uint8((word >> (i * 8)) & 0xFF)
		msg.WriteString(byteToString(b))
	}
	return msg.String()
}

func generalRegisters(conn net.Conn) {
	registersData := C.callGlobalRegistersCallback(globalRegistersCallback)
	registers := (*[1 << 28]C.uint)(unsafe.Pointer(registersData))[:38:38]
	var msg strings.Builder
	for i := 0; i < len(registers); i++ {
		register := wordToString(uint32(registers[i]))
		msg.WriteString(register)
	}
	for i := 38; i <= 72; i++ {
		msg.WriteString("xxxxxxxx")
	}
	send(conn, msg.String())
}

func reply(conn net.Conn, packet string) {
	split := strings.Split(packet, ":")
	method := split[0]
	switch method {
	case "qSupported":
		send(conn, "PacketSize=3fff;QStartNoAckMode+;")
	case "QStartNoAckMode":
		send(conn, "OK")
		ackDisabled = true
	case "?":
		send(conn, "S05")
	case "Hc-1":
		send(conn, "OK")
	case "Hg0":
		send(conn, "OK")
	case "qOffsets":
		send(conn, "Text=0;Data=0;Bss=0")
	case "g":
		generalRegisters(conn)
	default:
		send(conn, "")
	}
}
