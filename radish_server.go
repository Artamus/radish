package radish

import (
	"fmt"
	"net"
	"strconv"
)

type RadishServer struct {
	listener net.Listener
}

func NewRadishServer(port int) (*RadishServer, error) {
	l, err := net.Listen("tcp", "0.0.0.0:"+strconv.Itoa(port))
	if err != nil {
		return nil, fmt.Errorf("Failed to bind to port %d, %v", port, err)
	}

	return &RadishServer{
		listener: l,
	}, nil
}

func (r *RadishServer) Listen() {
	for {
		conn, err := r.listener.Accept()
		if err != nil {
			return
		}

		go handleConnection(conn)
	}
}

func (r *RadishServer) Close() {
	r.listener.Close()
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	for {
		buf := make([]byte, 128)
		conn.Read(buf)

		conn.Write([]byte("+PONG\r\n"))
	}
}
