package ratelimit

import (
	"math/rand"
	"time"
)

// NoOpRateLimiter is a simplistic version of RateLimiter that has no logic other than to
// blindly dump the messages sent down the ReceiveChan.
// This fails the Client tests, and allows for an example starting point for a real implementation.
type NoOpRateLimiter struct {
	reject  bool
	msgChan chan *Message
}

func NewNoOpRateLimiter(reject bool) *NoOpRateLimiter {
	return &NoOpRateLimiter{
		reject:  reject,
		msgChan: make(chan *Message),
	}
}

func (rl *NoOpRateLimiter) Start() error {
	return nil
}

func (rl *NoOpRateLimiter) Close() error {
	close(rl.msgChan)
	return nil
}

func (rl *NoOpRateLimiter) Send(msg *Message) error {
	if rand.Intn(100) < 10 { // reject ~10% of the messages
		if rl.reject {
			return ErrReject
		}
		time.Sleep(time.Millisecond * 100)
	}
	rl.msgChan <- msg
	return nil
}

func (rl *NoOpRateLimiter) ReceiveChan() <-chan *Message {
	return rl.msgChan
}
