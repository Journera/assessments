package ratelimit

import (
	"fmt"
	"github.com/brianvoe/gofakeit"
	"github.com/journera/assessments/common"
	"sync"
	"time"
)

type Client struct {
	limiter    RateLimiter
	msgCount   int
	msgsPerMin int
	randomMax  time.Duration

	done     sync.WaitGroup
	received common.LinkedList[*Message]
	failed   common.LinkedList[*Message]
}

func NewClient(limiter RateLimiter, msgCount, msgsPerMin int, randomMax time.Duration) *Client {
	c := &Client{
		limiter:    limiter,
		msgCount:   msgCount,
		msgsPerMin: msgsPerMin,
		randomMax:  randomMax,

		received: common.NewLinkedList[*Message](),
		failed:   common.NewLinkedList[*Message](),
	}
	c.done.Add(1)
	go c.receive()
	return c
}

func (c *Client) Run() {
	log.Info().Msgf("Running | Msgs:%d, PerMin:%d", c.msgCount, c.msgsPerMin)
	c.send()
	c.limiter.Close()
	log.Debug().Msg("Waiting for receive to finish")
	c.done.Wait()
	log.Info().Msgf("Complete | Received:%d, Failed:%d", c.received.Size(), c.failed.Size())
	c.Evaluate()
}

func (c *Client) send() {
	gofakeit.Seed(time.Now().Unix())
	log.Info().Msg("Starting sender")
	var throttle *time.Ticker
	if c.msgsPerMin > 0 {
		throttle = time.NewTicker(time.Minute / time.Duration(c.msgsPerMin))
		defer throttle.Stop()
	}

	var err error
	for i := 0; i < c.msgCount; i++ {
		msg := NewMessage(i,
			gofakeit.HackerAbbreviation(),
			fmt.Sprintf("%s %s %s", gofakeit.HackerVerb(), gofakeit.HackerAdjective(), gofakeit.HackerNoun()))
		log.Debug().Stringer("Msg", msg).Msg("Sending")
		err = c.limiter.Send(msg)
		if err != nil {
			c.failed.AddLast(msg)
		}
		if throttle != nil {
			<-throttle.C
		}
	}
	log.Debug().Msg("Sender complete")
}

func (c *Client) receive() {
	log.Info().Msg("Starting receiver")
	rcvChan := c.limiter.ReceiveChan()
	for msg := range rcvChan {
		msg.ReceiveTime = time.Now()
		log.Debug().Int("Msg", msg.Id).TimeDiff("Delay", msg.ReceiveTime, msg.SendTime).Msg("Received")
		c.received.AddLast(msg)
	}
	log.Debug().Msg("Receiver complete")
	c.done.Done()
}

func (c *Client) Evaluate() {
	log.Info().Msg("Evaluate")
}
