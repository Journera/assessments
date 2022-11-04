package ratelimit

import (
	"github.com/journera/assessments/common"
)

var (
	log = common.ProvideLog()
)

func ProvideRateLimiter(msgPerMin int, reject bool) RateLimiter {
	// TODO: replace with your actual RateLimiter
	return NewFakeRateLimiter(reject)
}
