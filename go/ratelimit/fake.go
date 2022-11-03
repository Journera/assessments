package ratelimit

import (
	"math/rand"
	"time"
)

// FakeRateLimiter is a simplistic version of RateLimiter that has no logic other than to
// blindly dump the messages sent down the ReceiveChan and reject some on a percentage basis.
// This fails the Client tests, and allows for an example starting point for a real implementation.
type FakeRateLimiter struct {
	reject  bool
	msgChan chan *Message
}

func NewFakeRateLimiter(reject bool) *FakeRateLimiter {
	return &FakeRateLimiter{
		reject:  reject,
		msgChan: make(chan *Message),
	}
}

func (rl *FakeRateLimiter) Start() error {
	return nil
}

func (rl *FakeRateLimiter) Close() error {
	close(rl.msgChan)
	return nil
}

func (rl *FakeRateLimiter) Send(msg *Message) error {
	if rand.Intn(100) < 10 { // reject ~10% of the messages
		if rl.reject {
			return ErrReject
		}
		time.Sleep(time.Millisecond * time.Duration(rand.Intn(50)+50))
	}
	rl.msgChan <- msg
	return nil
}

func (rl *FakeRateLimiter) ReceiveChan() <-chan *Message {
	return rl.msgChan
}
