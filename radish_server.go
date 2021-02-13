package radish

import (
	"fmt"
	"net"
	"strconv"
)

type RadishServer struct {
	listener net.Listener
	clients  []net.Conn
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

	serverChan := make(chan net.Conn)
	connChan := make(chan net.Conn)

	for {
		go func() {
			conn, err := r.listener.Accept()
			if err != nil {
				return
			}

			serverChan <- conn
		}()

		for _, conn := range r.clients {
			go func(conn net.Conn) {
				buf := make([]byte, 0)
				conn.Read(buf)
				connChan <- conn
			}(conn)
		}

		select {
		case newConn := <-serverChan:
			r.clients = append(r.clients, newConn)
		case existingConn := <-connChan:
			handleConnection(existingConn)
		}
	}
}

func (r *RadishServer) Close() {
	r.listener.Close()
}

func handleConnection(conn net.Conn) {
	buf := make([]byte, 128)
	conn.Read(buf)

	conn.Write([]byte("+PONG\r\n"))
}
