package ratelimit

import (
	"io"
)

// RateLimiter will limit the number of messages sent through it.
// Depending on configuration, it will either delay or reject messages to
// ensure a sender is not exceeding the limit.
type RateLimiter interface {
	io.Closer

	// Start the limiter (if needed) and initiate any internal activities. Does not block.
	Start() error

	// Send a message, if not over limit, will be delivered to ReceiveChan in near-real time.
	// If reject is enabled, may return an error.
	// If reject is disabled, may block until enough time has elapsed.
	Send(msg *Message) error

	// ReceiveChan return the channel that will receive messages that proceed through the limiter
	ReceiveChan() <-chan *Message
}
