package prat

import (
	"encoding/gob"
	"fmt"
	"net"
	"os"
	"time"
)

// MessageHandler is a callback which is invoked
// when the a new message is received
type MessageHandler func(Message)
type ErrorHandler func()

type Client struct {
	Conn            net.Conn
	Name            string
	messageHandlers []MessageHandler
	errorHandlers   []ErrorHandler
}

// Returns a Client connected to the given address
// The client will be listening to new messages from the host
func NewClient(address, name string) *Client {
	// connect to host
	conn, err := net.Dial("tcp", address)
	if err != nil {
		// could not connect to host
		// exit with error code
		fmt.Printf("Could not connect to %s\n", address)
		os.Exit(1)
	}

	client := Client{
		Conn:            conn,
		Name:            name,
		messageHandlers: make([]MessageHandler, 0),
		errorHandlers:   make([]ErrorHandler, 0),
	}

	go client.Listen(conn)
	return &client
}

func (c *Client) Listen(conn net.Conn) {
	defer conn.Close()
	for {
		var message Message
		dec := gob.NewDecoder(conn)
		err := dec.Decode(&message)
		if err != nil {
			// an error occured when decoding packet
			// exit with error and for the love of god
			// close the connection
			conn.Close()
			fmt.Println("Server error occured")
			for _, errorHandler := range c.errorHandlers {
				errorHandler()
			}
			break
		}
		for _, messageHandler := range c.messageHandlers {
			messageHandler(message)
		}
	}
}

// SendMessage creates a new Message with given
// string and sends it to the server
func (c *Client) SendMessage(message string) {
	// get time
	time := time.Now()
	// construct message struct
	msg := Message{
		Author: c.Name,
		Body:   message,
		Time:   time,
	}
	// encode and send to server
	enc := gob.NewEncoder(c.Conn)
	err := enc.Encode(msg)
	if err != nil {
		panic(err)
	}
}

func (c *Client) AddMessageHandler(messageHandler MessageHandler) {
	c.messageHandlers = append(c.messageHandlers, messageHandler)
}

func (c *Client) AddErrorHandler(errorHandler ErrorHandler) {
	c.errorHandlers = append(c.errorHandlers, errorHandler)
}
