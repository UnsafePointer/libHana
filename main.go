package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"net"
	"strings"
	"unicode"
)

type ScannerState byte

const (
	Beginning ScannerState = '$'
	End       ScannerState = '#'
)

var state = Beginning
var ackEnabled = true

func main() {
	lis, err := net.Listen("tcp", "localhost:1111")
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

func qSupported(conn net.Conn, params []string) {
	// send(conn, "PacketSize=3fff;QPassSignals+;QProgramSignals+;qXfer:libraries-svr4:read+;augmented-libraries-svr4-read+;qXfer:auxv:read+;qXfer:spu:read+;qXfer:spu:write+;qXfer:siginfo:read+;qXfer:siginfo:write+;qXfer:features:read+;QStartNoAckMode+;qXfer:osdata:read+;multiprocess+;fork-events+;vfork-events+;exec-events+;QDisableRandomization+;BreakpointCommands+;QAgent+;swbreak+;hwbreak+;vContSupported+;QThreadEvents+;no-resumed+")
	// send(conn, "PacketSize=3fff;qXfer:libraries-svr4;qXfer:features:read+;QStartNoAckMode+")
	// send(conn, "qSupported:multiprocess+;swbreak+;hwbreak+;qRelocInsn+;fork-events+;vfork-events+;exec-events+;vContSupported+;QThreadEvents+;no-resumed+")
	// send(conn, "PacketSize=3fff;QPassSignals+;QProgramSignals+;qXfer:libraries-svr4:read+;augmented-libraries-svr4-read+;qXfer:auxv:read+;qXfer:spu:read+;qXfer:spu:write+;qXfer:siginfo:read+;qXfer:siginfo:write+;qXfer:features:read+;QStartNoAckMode+;qXfer:osdata:read+;multiprocess+;fork-events+;vfork-events+;exec-events+;QNonStop+;QDisableRandomization+;qXfer:threads:read+;BreakpointCommands+;QAgent+;swbreak+;hwbreak+;qXfer:exec-file:read+;vContSupported+;QThreadEvents+;no-resumed+")
	// send(conn, "PacketSize=3fff")
	send(conn, "")
}

func vMustReplyEmpty(conn net.Conn, params []string) {
	send(conn, "")
}

func setThread(conn net.Conn, params string) {
	send(conn, "OK")
	// send(conn, "")
}

func queryTraceStatus(conn net.Conn) {
	// send(conn, "Trunning;tnotrun:0'")
	send(conn, "")
}

func stopReply(conn net.Conn) {
	// send(conn, "S 02'")
	send(conn, "S05")
}

func queryThreadInfo(conn net.Conn) {
	// send(conn, "m 1000'")
	send(conn, "")
}

func queryAdditionalThreadInfo(conn net.Conn) {
	// send(conn, "l'")
	send(conn, "")
}

func queryCurrentThreadID(conn net.Conn) {
	// send(conn, "QC 1000.1000'")
	send(conn, "")
}

func queryAttached(conn net.Conn) {
	// send(conn, "1")
	send(conn, "")
}

func queryRTOS(conn net.Conn) {
	send(conn, "")
}

func queryOffsets(conn net.Conn) {
	send(conn, "Text=0;Data=0;Bss=0")
}

func QStartNoAckMode(conn net.Conn) {
	ackEnabled = false
	send(conn, "OK")
}

func QProgramSignals(conn net.Conn) {
	send(conn, "OK")
}

func targetXML(conn net.Conn) {
	response := `l<?xml version="1.0"?><!DOCTYPE target SYSTEM "gdb-target.dtd"><target><architecture>mips</architecture><osabi>GNU/Linux</osabi><xi:include href="mips-cpu.xml"/><xi:include href="mips-cp0.xml"/><xi:include href="mips-fpu.xml"/><feature name="org.gnu.gdb.mips.linux"></feature></target>`
	send(conn, response)
}

func stripSpaces(str string) string {
	return strings.Map(func(r rune) rune {
		if unicode.IsSpace(r) {
			// if the character is a space, drop it
			return -1
		}
		// else keep it in the string
		return r
	}, str)
}

func mipsCPU(conn net.Conn) {
	dat, _ := ioutil.ReadFile("test.xml")
	str := string(dat)
	send(conn, stripSpaces(str))
}

func mipsCP0(conn net.Conn) {
	response := `l<?xml version="1.0"?><!DOCTYPE feature SYSTEM "gdb-target.dtd"><feature name="org.gnu.gdb.mips.cp0"><reg name="status" bitsize="32" regnum="32"/><reg name="badvaddr" bitsize="32" regnum="35"/><reg name="cause" bitsize="32" regnum="36"/></feature>`
	send(conn, response)
}

func mipsFPU(conn net.Conn) {
	response := `l<?xml version="1.0"?><!DOCTYPE feature SYSTEM "gdb-target.dtd"><feature name="org.gnu.gdb.mips.fpu"><reg name="f0" bitsize="32" type="ieee_single" regnum="38"/><reg name="f1" bitsize="32" type="ieee_single"/><reg name="f2" bitsize="32" type="ieee_single"/><reg name="f3" bitsize="32" type="ieee_single"/><reg name="f4" bitsize="32" type="ieee_single"/><reg name="f5" bitsize="32" type="ieee_single"/><reg name="f6" bitsize="32" type="ieee_single"/><reg name="f7" bitsize="32" type="ieee_single"/><reg name="f8" bitsize="32" type="ieee_single"/><reg name="f9" bitsize="32" type="ieee_single"/><reg name="f10" bitsize="32" type="ieee_single"/><reg name="f11" bitsize="32" type="ieee_single"/><reg name="f12" bitsize="32" type="ieee_single"/><reg name="f13" bitsize="32" type="ieee_single"/><reg name="f14" bitsize="32" type="ieee_single"/><reg name="f15" bitsize="32" type="ieee_single"/><reg name="f16" bitsize="32" type="ieee_single"/><reg name="f17" bitsize="32" type="ieee_single"/><reg name="f18" bitsize="32" type="ieee_single"/><reg name="f19" bitsize="32" type="ieee_single"/><reg name="f20" bitsize="32" type="ieee_single"/><reg name="f21" bitsize="32" type="ieee_single"/><reg name="f22" bitsize="32" type="ieee_single"/><reg name="f23" bitsize="32" type="ieee_single"/><reg name="f24" bitsize="32" type="ieee_single"/><reg name="f25" bitsize="32" type="ieee_single"/><reg name="f26" bitsize="32" type="ieee_single"/><reg name="f27" bitsize="32" type="ieee_single"/><reg name="f28" bitsize="32" type="ieee_single"/><reg name="f29" bitsize="32" type="ieee_single"/><reg name="f30" bitsize="32" type="ieee_single"/><reg name="f31" bitsize="32" type="ieee_single"/><reg name="fcsr" bitsize="32" group="float"/><reg name="fir" bitsize="32" group="float"/></feature>`
	send(conn, response)
}

func QNonStop(conn net.Conn) {
	send(conn, "")
}

func generalRegisters(conn net.Conn) {
	var msg strings.Builder
	for i := 0; i < 32; i++ {
		msg.WriteString("0000000000000000")
	}
	msg.WriteString("0000000000000000")
	msg.WriteString("0000000000000000")
	msg.WriteString("xxxxxxxxxxxxxxxx")
	msg.WriteString("xxxxxxxxxxxxxxxx")
	msg.WriteString("0000000000000000")
	// msg.WriteString("0c00000084b8f4f70c000000846f5b562089f1ff2889f1ff00a0f4f70000000001535b5682020000000000000000000000000000000000000000000000000000")
	// msg.WriteString("'")
	send(conn, msg.String())
	// send(conn, "E 00'")
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

func reply(conn net.Conn, packet string) {
	split := strings.Split(packet, ":")
	method := split[0]
	switch method {
	case "g":
		generalRegisters(conn)
	case "qSupported":
		qSupported(conn, split[1:])
	case "vMustReplyEmpty":
		vMustReplyEmpty(conn, split[1:])
	case "QStartNoAckMode":
		QStartNoAckMode(conn)
	case "QProgramSignals":
		QProgramSignals(conn)
	case "QNonStop":
		QNonStop(conn)
	case "m2ef0,4":
		send(conn, "")
	case "vFile":
		send(conn, "")
	case "m565b8271,1":
		send(conn, "")
	default:
		if strings.HasPrefix(method, "H") {
			setThread(conn, method[1:])
		} else if strings.HasPrefix(method, "q") {
			queryMethod := method[1:]
			switch queryMethod {
			case "Symbol":
				send(conn, "")
			case "TStatus":
				queryTraceStatus(conn)
			case "fThreadInfo":
				queryThreadInfo(conn)
			case "sThreadInfo":
				queryAdditionalThreadInfo(conn)
			case "C":
				queryCurrentThreadID(conn)
			case "Attached":
				queryAttached(conn)
			case "Offsets":
				queryOffsets(conn)
			case "Xfer":
				xfer := strings.Join(split[1:], ":")
				if strings.HasPrefix(xfer, "features:read:target.xml") {
					targetXML(conn)
				} else if strings.HasPrefix(xfer, "features:read:mips-cpu.xml") {
					mipsCPU(conn)
				} else if strings.HasPrefix(xfer, "features:read:mips-cp0.xml") {
					mipsCP0(conn)
				} else if strings.HasPrefix(xfer, "features:read:mips-fpu.xml") {
					mipsFPU(conn)
				} else {
					fmt.Println("Unhandled Xfer", xfer)
				}
			default:
				if strings.HasPrefix(queryMethod, "L") {
					queryRTOS(conn)
				} else {
					fmt.Println("Unhandled query method", queryMethod)
				}
			}
		} else if strings.HasPrefix(method, "?") {
			stopReply(conn)
		} else {
			fmt.Println("Unhandled ", method)
			send(conn, "")
		}
	}
}

func ack(conn net.Conn) {
	if ackEnabled {
		conn.Write([]byte("+"))
	}
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
			// fmt.Println("Empty received")
		case "+":
			// fmt.Println("ACK received")
			// ack(conn) ?
		case "-":
			// fmt.Println("NACK received")
			// ack(conn) ?
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
