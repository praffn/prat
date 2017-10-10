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
	Conn net.Conn
	Name string
}

// Returns a Client connected to the given address
// The client will be listening to new messages from the host
func NewClient(address, name string, messageHandler MessageHandler) *Client {
	// connect to host
	conn, err := net.Dial("tcp", address)
	if err != nil {
		panic(err)
	}

	// start go routine listening for messages
	go func() {
		for {
			var message Message
			// decode received message
			dec := gob.NewDecoder(conn)
			err := dec.Decode(&message)
			if err != nil {
				panic(err)
			}
			// invoke messagehandler with decoded message
			messageHandler(message)
		}
	}()

	// return new client struct
	return &Client{
		Conn: conn,
		Name: name,
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
