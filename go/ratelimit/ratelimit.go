package ratelimit

import (
	"io"
)

// RateLimiter will limit the number of messages sent through it.
// Depending on configuration, it will either delay or reject messages to
// ensure a sender is not exceeding the limit.
type RateLimiter interface {
	io.Closer

	// Start the limiter and initiate any internal activities
	Start() error

	// Send a message and check if the Sender has exceeded the rate
	Send(msg *Message) error

	// ReceiveChan return the channel that will receive messages that proceed through the limiter
	ReceiveChan() <-chan *Message
}
