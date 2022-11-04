package ratelimit

import (
	"github.com/journera/assessments/common"
	"time"
)

type Collector struct {
	limiter RateLimiter

	received *common.LinkedList[*Message]
}

func NewCollector(limiter RateLimiter) *Collector {
	c := &Collector{
		limiter: limiter,

		received: common.NewLinkedList[*Message](),
	}
	return c
}

// Run will collect all messages from the rate limiter.
// This will block until the limiter closes the channel
func (c *Collector) Run() {
	log.Info().Msg("Starting collector")
	rcvChan := c.limiter.ReceiveChan()
	for msg := range rcvChan {
		msg.ReceiveTime = time.Now()
		log.Trace().Int("Msg", msg.Id).TimeDiff("Delay", msg.ReceiveTime, msg.SendTime).Msg("Received")
		c.received.AddLast(msg)
	}
	log.Debug().Msg("Collector complete")
}
