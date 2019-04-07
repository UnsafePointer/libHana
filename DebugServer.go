// package name: libHana
package main

/*
#cgo CFLAGS:
typedef unsigned int* (*GlobalRegistersCallbackType)();
unsigned int* callGlobalRegistersCallback(GlobalRegistersCallbackType callback);
typedef unsigned char* (*MemoryReadCallbackType)(unsigned int address, unsigned int length);
unsigned char* callMemoryReadCallback(MemoryReadCallbackType callback, unsigned int address, unsigned int length);
typedef void (*AddBreakpointCallbackType)(unsigned int address);
void callAddBreakpointCallback(AddBreakpointCallbackType callback, unsigned int address);
typedef void (*ContinueCallbackType)();
void callContinueCallback(ContinueCallbackType callback);
*/
import "C"

import (
	"bufio"
	"bytes"
	"fmt"
	"net"
	"strconv"
	"strings"
	"unsafe"
)

type ScannerState byte

const (
	Beginning ScannerState = '$'
	End       ScannerState = '#'
)

var (
	state = Beginning
	done  chan struct{}
)

var globalRegistersCallback C.GlobalRegistersCallbackType
var readMemoryCallback C.MemoryReadCallbackType
var addBreakpointCallback C.AddBreakpointCallbackType
var continueCallback C.ContinueCallbackType
var ackDisabled = false

//export SetGlobalRegistersCallback
func SetGlobalRegistersCallback(fn C.GlobalRegistersCallbackType) {
	globalRegistersCallback = fn
}

//export SetReadMemoryCallback
func SetReadMemoryCallback(fn C.MemoryReadCallbackType) {
	readMemoryCallback = fn
}

//export SetAddBreakpointCallback
func SetAddBreakpointCallback(fn C.AddBreakpointCallbackType) {
	addBreakpointCallback = fn
}

//export SetContinueCallback
func SetContinueCallback(fn C.ContinueCallbackType) {
	continueCallback = fn
}

//export NotifyStopped
func NotifyStopped() {
	close(done)
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
	go handleConnection(conn)
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
	for {
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

func memoryRead(address uint32, length uint32, conn net.Conn) {
	memoryData := C.callMemoryReadCallback(readMemoryCallback, C.uint(address), C.uint(length))
	memory := (*[1 << 28]C.uint)(unsafe.Pointer(memoryData))[:length:length]
	var msg strings.Builder
	for i := uint32(0); i < length; i++ {
		value := byteToString(uint8(memory[i]))
		msg.WriteString(value)
	}
	send(conn, msg.String())
}

func addBreakpoint(address uint32) {
	C.callAddBreakpointCallback(addBreakpointCallback, C.uint(address))
}

func addLoadWatchpoint(address uint32) {

}

func addStoreWatchpoint(address uint32) {

}

func continueProgram(conn net.Conn) {
	C.callContinueCallback(continueCallback)
	done = make(chan struct{})
	<-done
	send(conn, "S00")
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
		if strings.HasPrefix(method, "c") {
			continueProgram(conn)
		} else if strings.HasPrefix(method, "m") {
			params := strings.Split(method[1:], ",")
			address, _ := strconv.ParseUint(params[0], 10, 32)
			length, _ := strconv.ParseUint(params[1], 10, 32)
			memoryRead(uint32(address), uint32(length), conn)
		} else if strings.HasPrefix(method, "Z") {
			params := strings.Split(method[1:], ",")
			breakType, _ := strconv.ParseUint(params[0], 10, 32)
			address, _ := strconv.ParseUint(params[1], 16, 32)
			kind, _ := strconv.ParseUint(params[2], 10, 32)
			if kind != 4 {
				// libHana only supports MIPS 32-Bits
				// https://sourceware.org/gdb/onlinedocs/gdb/MIPS-Breakpoint-Kinds.html
				send(conn, "E00")
			}
			switch breakType {
			case 0:
				addBreakpoint(uint32(address))
				send(conn, "OK")
			case 2:
				addLoadWatchpoint(uint32(address))
				send(conn, "OK")
			case 3:
				addStoreWatchpoint(uint32(address))
				send(conn, "OK")
			default:
				send(conn, "")
			}
		} else {
			send(conn, "")
		}
	}
}
