package radish

import (
	"fmt"
	"io"
	"net"
)

type client struct {
	socket net.Conn
	buffer []byte
}

func newClient(socket net.Conn) *client {
	return &client{socket: socket, buffer: make([]byte, 0)}
}

func (c *client) readAvailable() error {
	buf := make([]byte, 1024)
	_, err := c.socket.Read(buf)
	if err != nil && err == io.EOF {
		return fmt.Errorf("connection appears to have closed, %v", err)
	}

	c.buffer = append(c.buffer, buf...)
	return nil
}

func (c *client) write(message string) {
	c.socket.Write([]byte(message))
}
