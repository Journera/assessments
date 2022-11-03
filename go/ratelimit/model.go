package ratelimit

import (
	"errors"
	"fmt"
	"github.com/journera/assessments/common"
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
	return fmt.Sprintf("[%d] %s", m.Id, m.Text)
}

type Stats struct {
	Sender      string
	Messages    *common.LinkedList[*Message]
	Missing     []int
	MinTime     time.Duration
	MaxTime     time.Duration
	AverageTime time.Duration
	TotalTime   time.Duration
}

func NewStats(msg *Message) *Stats {
	return &Stats{
		Sender:   msg.Sender,
		Messages: common.NewLinkedList(msg),
		Missing:  make([]int, 0, 8),
	}
}
