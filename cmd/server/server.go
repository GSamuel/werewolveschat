package main

import (
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
	connections []*chat.Connection
	data        map[string]string
)

func SocketServer(port int) {

	listen, err := net.Listen("tcp4", ":"+strconv.Itoa(port))
	defer listen.Close()
	if err != nil {
		log.Fatalf("Socket listen port %d failed,%s", port, err)
		os.Exit(1)
	}
	log.Printf("Begin listen port: %d", port)

	data = make(map[string]string)

	go ProccessInput()

	for {
		conn, err := listen.Accept()
		if err != nil {
			log.Fatalln(err)
			continue
		}
		c := chat.NewConnection(conn)
		connections = append(connections, c)
		c.Start()
	}

}

func ProccessInput() {
	for {
		for i, conn1 := range connections {
			if conn1.Running() {
				select {
				case msg := <-conn1.Output:
					ProcessCommand(msg, i)
				default:
				}
			}
		}
	}
}

func ProcessCommand(msg string, i int) {
	nameKey := fmt.Sprintf("%v:name", i)
	name := data[nameKey]

	if strings.Index(msg, "/name ") == 0 {
		newName := strings.TrimPrefix(msg, "/name ")
		data[nameKey] = newName
		if name == "" {
			SendToAllExcept(fmt.Sprintf("%s (%v) joined the room", newName, i), i)
		} else {
			SendToAllExcept(fmt.Sprintf("%s (%v) changed name to %s", name, i, newName), i)
		}
		return
	}

	if strings.Index(msg, "/help") == 0 {
		//msg = strings.TrimPrefix(msg, "/help ")
		SendToOne("\n\t/help\t\tlists all commands\n\t/name [name]\t change your name", i)
		return
	}

	if strings.Index(msg, "/list") == 0 {
		result := ""
		for j := 0; j < len(connections); j++ {
			if connections[j].Running() {
				result = fmt.Sprintf("%s\t%s (%v)", result, data[fmt.Sprintf("%v:name", j)], j)
			}
		}
		SendToOne(result, i)
		return
	}

	SendToAllExcept(fmt.Sprintf("%s (%v): %s", name, i, msg), i)
}

func SendToOne(msg string, n int) {
	conn := connections[n]
	if conn.Running() {
		conn.Writer().Write([]byte(msg))
		conn.Writer().Flush()
	}
}

func SendToAllExcept(msg string, n int) {
	for i, conn2 := range connections {
		if conn2.Running() && i != n {
			conn2.Writer().Write([]byte(msg))
			conn2.Writer().Flush()
		}
	}
}

func initLog() {
	f, err := os.OpenFile("server.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
		defer f.Close()
	}

	log.SetOutput(f)
}

func isTransportOver(data string) (over bool) {
	over = strings.HasSuffix(data, "\r\n\r\n")
	return
}

func main() {
	port := 10001

	//initLog()
	SocketServer(port)
}
