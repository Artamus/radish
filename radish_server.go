package radish

import (
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"
	"time"
)

// Server encapsulates a TCP server that can serve RESP commands
type Server struct {
	listener      net.Listener
	connClients   map[net.Conn]*client
	storage       map[string]string
	cancellations map[string]chan struct{}
}

// NewServer is used as the constructor to assemble new servers
func NewServer(port int, keyValueStorage map[string]string) (*Server, error) {
	l, err := net.Listen("tcp", "0.0.0.0:"+strconv.Itoa(port))
	if err != nil {
		return nil, fmt.Errorf("Failed to bind to port %d, %v", port, err)
	}

	return &Server{
		listener:      l,
		connClients:   make(map[net.Conn]*client),
		storage:       keyValueStorage,
		cancellations: make(map[string]chan struct{}),
	}, nil
}

// Listen makes the server start listening for incoming connections on the port specified in the constructor
func (r *Server) Listen() {

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

		for conn := range r.connClients {
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

// Close stops the listener from accepting any more incoming connections
func (r *Server) Close() {
	r.listener.Close()
}

func (r *Server) handleClient(client *client) {
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

func (r *Server) handleCommand(client *client, command *command) {

	switch strings.ToUpper(command.action) {
	case "PING":
		client.write("+PONG\r\n")
	case "ECHO":
		firstArg := command.args[0]
		response := fmt.Sprintf("+%s\r\n", firstArg)
		client.write(response)
	case "GET":
		key := command.args[0]
		value, ok := r.storage[key]
		if !ok {
			client.write("$-1\r\n")
			return
		}
		client.write(fmt.Sprintf("+%s\r\n", value))
	case "SET":
		key := command.args[0]
		value := command.args[1]

		cancelExpiry(r.cancellations, key)

		if len(command.args) > 2 {

			timeUnit := time.Millisecond
			if strings.ToUpper(command.args[2]) == "EX" {
				timeUnit = time.Second
			}
			amount, _ := strconv.Atoi(command.args[3])

			cancellation := make(chan struct{})
			r.cancellations[key] = cancellation
			go r.expireKey(key, time.Duration(amount)*timeUnit, cancellation)
		}

		r.storage[key] = value

		client.write("+OK\r\n")
	default:
		log.Printf("unknown command '%s'", command.action)
		response := fmt.Sprintf("-ERR unknown command '%s'\r\n", command.action)
		client.write(response)
	}
}

func (r *Server) expireKey(key string, after time.Duration, cancel chan struct{}) {

	select {
	case <-time.After(after):
		delete(r.storage, key)
	case <-cancel:
	}

	delete(r.cancellations, key)
}

func cancelExpiry(cancellations map[string]chan struct{}, key string) {
	cancel, ok := cancellations[key]
	if ok {
		cancel <- struct{}{}
	}
}
