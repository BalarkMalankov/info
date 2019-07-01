package main

import (
	"fmt"
	"net"
	"bufio"
	"encoding/json"
	"strings"
	"strconv"
)

var lib map[int]string = make(map[int]string)

func hostname(name string, port int) string {
	return fmt.Sprintf("%s:%d", name, port)
}

func readMsg(con net.Conn) string {
	r := bufio.NewReader(con)
	msg,_ := r.ReadString('\n')
	return strings.TrimSpace(msg)
}

func send(port int, name, msg string) string {
	con, _ := net.Dial("tcp", hostname(name,port))
	fmt.Fprintln(con, msg)
	return readMsg(con)
}

func sendNR(port int, name, msg string) {
	con, _ := net.Dial("tcp", hostname(name,port))
	fmt.Fprintln(con, msg)
}

func cliAdder(port int) {
	for dir, name := range lib {
		sendNR(dir + 2, name,fmt.Sprintf("%d",port)) // puerto + 2
	}
}

func handleRegister(con net.Conn) {
	defer con.Close()
	msg := readMsg(con)
	port,_ := strconv.Atoi(msg)
	cliAdder(port)
	jlib, _ := json.Marshal(lib)
	fmt.Fprintln(con, string(jlib))
	lib[port] = "localhost" // seccion crítica
	fmt.Println(lib)
}

func servRegister(portbase int) {
	ln,_ := net.Listen("tcp", hostname("localhost", portbase+1))
	defer ln.Close()
	for {
		con,_ := ln.Accept()
		go handleRegister(con)
	}
}

func cliRegister(name string, servport int, myport int) {
	resp := send(servport + 1, name, fmt.Sprintf("%d", myport))
	temp := make(map[int]string)
	_ = json.Unmarshal([]byte(resp), &temp)
	for port, na := range temp {
		lib[port] = na
	}
	fmt.Println(lib)
}

func handleAdder(con net.Conn) {
	defer con.Close()
	port,_ := strconv.Atoi(readMsg(con))
	lib[port] = "localhost" //sección crítica
	fmt.Println(lib)
}

func servAdder(name string, portbase int) {
	ln,_ := net.Listen("tcp", hostname(name, portbase + 2))
	defer ln.Close()
	for {
		cn,_ := ln.Accept()
		go handleAdder(cn)
	}

}

func main() {
	name := "localhost"
	port := 0
	fmt.Scanf("%d\n", &port)
	go servRegister(port)
	friendPort := 0
	fmt.Scanf("%d\n", &friendPort)
	if port != friendPort {
		lib[friendPort] = name
		cliRegister(name, friendPort, port)
	}
	servAdder(name, port)
}