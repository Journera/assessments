package ratelimit

import (
	"errors"
	"fmt"
	"time"
)

var (
	ErrReject = errors.New("too many messages")
)

type Message struct {
	Id          int
	Sender      string
	SendTime    time.Time
	ReceiveTime time.Time
	Text        string
}

func NewMessage(id int, sender, text string) *Message {
	return &Message{
		Id:       id,
		Sender:   sender,
		SendTime: time.Now(),
		Text:     text,
	}
}

func (m *Message) String() string {
	return fmt.Sprintf("%d %s [%s]", m.Id, m.Sender, m.Text)
}
