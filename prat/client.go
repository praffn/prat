package prat

import (
	"encoding/gob"
	"net"
	"time"
)

// MessageHandler is a callback which is invoked
// when the a new message is received
type MessageHandler func(Message)

type Client struct {
	Conn            net.Conn
	Name            string
	messageHandlers []MessageHandler
}

// Returns a Client connected to the given address
// The client will be listening to new messages from the host
func NewClient(address, name string) *Client {
	// connect to host
	conn, err := net.Dial("tcp", address)
	if err != nil {
		panic(err)
	}

	client := Client{
		Conn:            conn,
		Name:            name,
		messageHandlers: make([]MessageHandler, 0),
	}

	go client.Listen(conn)
	return &client
}

func (c *Client) Listen(conn net.Conn) {
	for {
		var message Message
		dec := gob.NewDecoder(conn)
		err := dec.Decode(&message)
		if err != nil {
			panic(err)
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
