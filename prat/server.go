package prat

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"io"
	"log"
	"net"
)

type Server struct {
	Connections map[net.Conn]bool
	NewConCh    chan net.Conn
	DeadConCh   chan net.Conn
	Messages    chan Message
	Logger      *log.Logger
}

func NewServer(logger *log.Logger) *Server {
	return &Server{
		Connections: make(map[net.Conn]bool),
		NewConCh:    make(chan net.Conn),
		DeadConCh:   make(chan net.Conn),
		Messages:    make(chan Message),
		Logger:      logger,
	}
}

func (s *Server) Start(port int) {
	// create address from port (like ":port")
	address := fmt.Sprintf(":%d", port)
	// start tcp listener on address
	listener, err := net.Listen("tcp", address)
	s.Logger.Printf("Server started, listening on port %d\n", port)
	if err != nil {
		s.Logger.Fatal(err)
		panic(err)
	}

	// start go routine accepting connection
	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				s.Logger.Fatal(err)
				panic(err)
			}
			// send connection to new connection channel
			s.NewConCh <- conn
		}
	}()

	// do forever
	for {
		select {
		case message := <-s.Messages:
			// on new message received
			s.OnNewMessage(message)
		case conn := <-s.NewConCh:
			// on new connection received
			s.OnNewConnection(conn)
		case conn := <-s.DeadConCh:
			// on dead connection
			s.OnDeadConnection(conn)
		}
	}
}

func (s *Server) OnNewMessage(message Message) {
	// create empty byte buffer
	var buf bytes.Buffer
	// create gob encoder to buffer
	enc := gob.NewEncoder(&buf)
	// encode message
	err := enc.Encode(message)
	if err != nil {
		s.Logger.Fatal(err)
	}

	// loop through connections
	for conn := range s.Connections {
		s.Logger.Printf("Sending message to %s\n", conn.RemoteAddr().String())
		// write to connection encoded message
		conn.Write(buf.Bytes())
		if err != nil {
			s.Logger.Printf("Error sending message to %s\n", conn.RemoteAddr().String())
			// if error while writing to connection
			// send connection to dead connection channel
			s.DeadConCh <- conn
		}
	}
}

func (s *Server) OnNewConnection(conn net.Conn) {
	s.Logger.Printf("New connection from %s\n", conn.RemoteAddr().String())
	// add to connection map received connection
	s.Connections[conn] = true
	// start go routine listening on connection
	go s.ConnectionListener(conn)
}

func (s *Server) ConnectionListener(conn net.Conn) {
	// forever
	for {
		// create decoder on connection
		dec := gob.NewDecoder(conn)
		// create empty message
		var message Message
		// decode received message
		err := dec.Decode(&message)
		if err == io.EOF {
			s.Logger.Printf("Connection from %s EOF'ed\n", conn.RemoteAddr().String())
			// if eof send connection to dead connection channel
			s.DeadConCh <- conn
			break
		} else if err != nil {
			panic(err)
		} else {
			s.Logger.Printf("Received message from %s: \"%s\"\n", conn.RemoteAddr().String(), message.Body)
			// send decoded message to message channel
			s.Messages <- message
		}
	}
}

func (s *Server) OnDeadConnection(conn net.Conn) {
	s.Logger.Printf("Connection from %s disconnected\n", conn.RemoteAddr().String())
	// remove connection from connections map
	delete(s.Connections, conn)
}
