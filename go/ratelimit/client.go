package ratelimit

import (
	"fmt"
	"github.com/brianvoe/gofakeit"
	"github.com/journera/assessments/common"
	"time"
)

type Client struct {
	limiter    RateLimiter
	sender     string
	msgCount   int
	msgsPerMin int
	failed     *common.LinkedList[*Message]
}

func NewClient(limiter RateLimiter, sender string, msgCount, msgsPerMin int) *Client {
	c := &Client{
		limiter:    limiter,
		sender:     sender,
		msgCount:   msgCount,
		msgsPerMin: msgsPerMin,
		failed:     common.NewLinkedList[*Message](),
	}
	return c
}

// Run will send all messages and block until complete
func (c *Client) Run() {
	log.Debug().
		Str("Sender", c.sender).
		Int("Msgs", c.msgCount).
		Int("PerMin", c.msgsPerMin).
		Msg("Starting Client")

	var throttle *time.Ticker
	if c.msgsPerMin > 0 {
		throttle = time.NewTicker(time.Minute / time.Duration(c.msgsPerMin))
		defer throttle.Stop()
	}

	var err error
	for i := 0; i < c.msgCount; i++ {
		msg := NewMessage(i, c.sender,
			fmt.Sprintf("%s %s %s", gofakeit.HackerVerb(), gofakeit.HackerAdjective(), gofakeit.HackerNoun()))
		log.Trace().Str("Sender", c.sender).Stringer("Msg", msg).Msg("Sending")
		err = c.limiter.Send(msg)
		if err != nil {
			c.failed.AddLast(msg)
		}
		if throttle != nil {
			<-throttle.C
		}
	}
	log.Debug().Str("Sender", c.sender).Msg("Client complete")
}
