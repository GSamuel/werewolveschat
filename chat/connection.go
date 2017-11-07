package chat

import (
	"bufio"
	"io"
	"log"
	"net"
	"strings"
)

type Connection struct {
	net.Conn
	r       *bufio.Reader
	w       *bufio.Writer
	Output  chan string
	started bool
	running bool
}

func (c *Connection) run() {

	defer c.Close()

	var (
		buf = make([]byte, 1024)
	)

	for {
		n, err := c.r.Read(buf)
		data := string(buf[:n])

		switch err {
		case io.EOF:
			c.running = false
			return
		case nil:
			c.Output <- data
			if c.isTransportOver(data) {
				c.running = false
				return
			}
		default:
			c.running = false
			log.Println("Receive data failed: %s", err)
			//log.Fatalf("Receive data failed:%s", err)
			return
		}

	}
}

func (c *Connection) Writer() *bufio.Writer {
	return c.w
}

func (c *Connection) Start() {
	if !c.running {
		c.started = true
		c.running = true
		go c.run()
	}
}

func (c *Connection) Started() bool {
	return c.started
}

func (c *Connection) Running() bool {
	return c.running
}

func (c *Connection) isTransportOver(data string) (over bool) {
	over = strings.HasSuffix(data, "\r\n\r\n")
	return
}

func NewConnection(conn net.Conn) *Connection {
	r := bufio.NewReader(conn)
	w := bufio.NewWriter(conn)
	return &Connection{conn, r, w, make(chan string, 10), false, false}
}
