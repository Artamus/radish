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

func (c *client) consumeCommand() *command {
	decoded, err := Decode(string(c.buffer))
	if err != nil && err == ErrIncompleteRESP {
		return nil
	}

	if len(decoded) == 0 {
		return nil
	}

	c.buffer = make([]byte, 0)

	return newCommand(decoded[0], decoded[1:])
}

func (c *client) write(message string) {
	c.socket.Write([]byte(message))
}

type command struct {
	action string
	args   []string
}

func newCommand(action string, args []string) *command {
	return &command{action: action, args: args}
}
