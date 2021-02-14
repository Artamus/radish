package radish

import "net"

type client struct {
	socket net.Conn
	buffer []byte
}

func newClient(socket net.Conn) *client {
	return &client{socket: socket, buffer: make([]byte, 0)}
}

func (c *client) readAvailable() {
	buf := make([]byte, 1024)
	c.socket.Read(buf)

	c.buffer = append(c.buffer, buf...)
}

func (c *client) write(message string) {
	c.socket.Write([]byte(message))
}
