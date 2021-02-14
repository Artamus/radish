package radish

import (
	"fmt"
	"net"
	"strconv"
	"strings"
)

type RadishServer struct {
	listener    net.Listener
	connClients map[net.Conn]*client
	storage     map[string]string
}

func NewRadishServer(port int, keyValueStorage map[string]string) (*RadishServer, error) {
	l, err := net.Listen("tcp", "0.0.0.0:"+strconv.Itoa(port))
	if err != nil {
		return nil, fmt.Errorf("Failed to bind to port %d, %v", port, err)
	}

	return &RadishServer{
		listener:    l,
		connClients: make(map[net.Conn]*client),
		storage:     keyValueStorage,
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

			client, ok := r.connClients[existingConn]
			if !ok {
				// TODO: Sometimes there are messages from a channel with the pointer address 0x0, not sure why, but they are not my clients
				break
			}
			r.handleClient(client)
		}
	}
}

func (r *RadishServer) Close() {
	r.listener.Close()
}

func (r *RadishServer) handleClient(client *client) {
	err := client.readAvailable()
	if err != nil {
		delete(r.connClients, client.socket)
		return
	}

	for {
		command := client.consumeCommand()
		if command == nil {
			break
		}

		r.handleCommand(client, command)
	}
}

func (r *RadishServer) handleCommand(client *client, command *command) {

	switch strings.ToUpper(command.action) {
	case "PING":
		client.write("+PONG\r\n")
	case "ECHO":
		firstArg := command.args[0]
		response := fmt.Sprintf("+%s\r\n", firstArg)
		client.write(response)
	case "GET":
		key := command.args[0]
		keyString, ok := key.(string)
		if !ok {
			client.write(fmt.Sprintf("-ERR unknown key '%v'\r\n", key))
			return
		}
		value, ok := r.storage[keyString]
		if !ok {
			client.write("$-1\r\n")
			return
		}
		client.write(fmt.Sprintf("+%s\r\n", value))
	default:
		response := fmt.Sprintf("-ERR unknown command '%s'\r\n", command.action)
		client.write(response)
	}
}
