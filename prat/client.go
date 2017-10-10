package prat

import (
	"encoding/gob"
	"net"
	"time"
)

type MessageHandler func(Message)

type Client struct {
	Conn           net.Conn
	Name           string
	MessageHandler MessageHandler
}

func NewClient(address, name string, messageHandler MessageHandler) *Client {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		panic(err)
	}

	go func() {
		for {
			var message Message
			dec := gob.NewDecoder(conn)
			err := dec.Decode(&message)
			if err != nil {
				panic(err)
			}
			messageHandler(message)
		}
	}()

	return &Client{
		Conn:           conn,
		Name:           name,
		MessageHandler: messageHandler,
	}
}

func (c *Client) SendMessage(message string) {
	time := time.Now()
	msg := Message{
		Author: c.Name,
		Body:   message,
		Time:   time,
	}
	enc := gob.NewEncoder(c.Conn)
	err := enc.Encode(msg)
	if err != nil {
		panic(err)
	}
}
