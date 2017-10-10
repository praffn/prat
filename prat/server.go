package prat

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"time"
)

type Server struct {
	Connections map[net.Conn]bool
	NewConCh    chan net.Conn
	DeadConCh   chan net.Conn
	Messages    chan Message
	Logger      *log.Logger
}

func NewServer() *Server {
	return &Server{
		Connections: make(map[net.Conn]bool),
		NewConCh:    make(chan net.Conn),
		DeadConCh:   make(chan net.Conn),
		Messages:    make(chan Message),
		Logger:      log.New(os.Stdout, "", 0),
	}
}

func NewServerWithLogger(logger *log.Logger) *Server {
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
	s.Log(fmt.Sprintf("Server started, listening on port %d", port))
	if err != nil {
		s.Fatal(err)
		panic(err)
	}

	// start go routine accepting connection
	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				s.Fatal(err)
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
		s.Fatal(err)
	}

	// loop through connections
	for conn := range s.Connections {
		s.Log(fmt.Sprintf("Sending message to %s", conn.RemoteAddr().String()))
		// write to connection encoded message
		conn.Write(buf.Bytes())
		if err != nil {
			s.Log(fmt.Sprintf("Error sending message to %s", conn.RemoteAddr().String()))
			// if error while writing to connection
			// send connection to dead connection channel
			s.DeadConCh <- conn
		}
	}
}

func (s *Server) OnNewConnection(conn net.Conn) {
	s.Log(fmt.Sprintf("New connection from %s", conn.RemoteAddr().String()))
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
			s.Log(fmt.Sprintf("Connection from %s EOF'ed", conn.RemoteAddr().String()))
			// if eof send connection to dead connection channel
			s.DeadConCh <- conn
			break
		} else if err != nil {
			panic(err)
		} else {
			s.Log(fmt.Sprintf("Received message from %s: \"%s\"", conn.RemoteAddr().String(), message.Body))
			// send decoded message to message channel
			s.Messages <- message
		}
	}
}

func (s *Server) OnDeadConnection(conn net.Conn) {
	s.Log(fmt.Sprintf("Connection from %s disconnected", conn.RemoteAddr().String()))
	// remove connection from connections map
	delete(s.Connections, conn)
}

func (s *Server) Log(message string) {
	s.Logger.Printf("[%s] %s\n", time.Now().Format("2006-01-02 15:04:05"), message)
}

func (s *Server) Fatal(err error) {
	s.Logger.Fatalf("[%s] <ERROR> %v\n", time.Now().Format("2006-01-02 15:04:05"), err)
}
