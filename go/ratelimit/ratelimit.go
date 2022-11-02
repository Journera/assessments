package ratelimit

import (
	"io"
)

type RateLimiter interface {
	io.Closer

	// Start the limiter and initiate any internal activities
	Start() error

	// Send a message and check if the Sender has exceeded the rate
	Send(msg *Message) error

	// ReceiveChan return the channel that will receive messages that proceed through the limiter
	ReceiveChan() <-chan *Message
}
