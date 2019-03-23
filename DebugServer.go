// package name: libHana
package main

/*
#cgo CFLAGS:
typedef int* (*GlobalRegistersCallbackType)();
int* callGlobalRegistersCallback(GlobalRegistersCallbackType callback);
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

type CallbackType func(C.int)

type ScannerState byte

const (
	Beginning ScannerState = '$'
	End       ScannerState = '#'
)

var (
	state = Beginning
)

var globalRegistersCallback C.GlobalRegistersCallbackType

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

func generalRegisters(conn net.Conn) {
	registersData := C.callGlobalRegistersCallback(globalRegistersCallback)
	registers := (*[1 << 28]C.int)(unsafe.Pointer(registersData))[:3:3]
	// TODO:
	fmt.Println("Registers: ", registers)
}

func reply(conn net.Conn, packet string) {
	split := strings.Split(packet, ":")
	method := split[0]
	switch method {
	case "g":
		fmt.Println("General registers")
		generalRegisters(conn)
	default:
		send(conn, "")
	}
}
