package radish

import (
	"fmt"
	"net"
	"strconv"
)

type RadishServer struct {
	listener    net.Listener
	connClients map[net.Conn]*client
}

func NewRadishServer(port int) (*RadishServer, error) {
	l, err := net.Listen("tcp", "0.0.0.0:"+strconv.Itoa(port))
	if err != nil {
		return nil, fmt.Errorf("Failed to bind to port %d, %v", port, err)
	}

	return &RadishServer{
		listener:    l,
		connClients: make(map[net.Conn]*client),
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

		for conn, _ := range r.connClients {
			go func(conn net.Conn) {
				buf := make([]byte, 0)
				conn.Read(buf)
				connChan <- conn
			}(conn)
		}

		select {
		case newConn := <-serverChan:
			r.connClients[newConn] = newClient(newConn)
		case existingConn := <-connChan:
			client := r.connClients[existingConn]
			handleClient(client)
		}
	}
}

func (r *RadishServer) Close() {
	r.listener.Close()
}

func handleClient(client *client) {
	client.readAvailable()
	client.write("+PONG\r\n")
}
