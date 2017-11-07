package main

import (
	"bufio"
	"fmt"
	"github.com/GSamuel/werewolveschat/chat"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
)

const (
	StopCharacter = "\r\n\r\n"
)

var (
	name = ""
)

func SocketClient(ip string, port int) {
	addr := strings.Join([]string{ip, strconv.Itoa(port)}, ":")
	conn, err := net.Dial("tcp", addr)

	if err != nil {
		log.Println(err)
		return
	}

	defer conn.Close()

	fmt.Println("Connected to server")

	conn.Write([]byte("/name " + name)) //set name on server

	c := chat.NewConnection(conn)
	c.Start()

	go keepReading(c)

	for c.Running() {
		message := ReadInput()
		conn.Write([]byte(message))
	}

}

func ReadInput() string {
	reader := bufio.NewReader(os.Stdin)
	//fmt.Print("You: ")
	message, _ := reader.ReadString('\n')
	message = strings.TrimSpace(message)
	return message
}

func keepReading(conn *chat.Connection) {
	for conn.Running() {
		msg := <-conn.Output
		fmt.Println(msg)
	}
}

func main() {

	var (
		ips   = []string{"127.0.0.1", "31.200.213.152"}
		ports = []int{10001, 10001}
	)

	fmt.Print("Enter name: ")
	name = ReadInput()
	for {
		for i, ip := range ips {
			port := ports[i]
			fmt.Println("Connection...")
			SocketClient(ip, port)
			chat.Clear()
		}
	}
}
