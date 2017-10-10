package prat

import "time"

type Message struct {
	Author string
	Body   string
	Time   time.Time
}

func NewMessage(author, body string, time time.Time) Message {
	return Message{author, body, time}
}
