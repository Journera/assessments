package ratelimit

import (
	"github.com/journera/assessments/common"
	"github.com/spf13/cobra"
	"time"
)

var (
	log = common.ProvideLog()
)

func ProvideRateLimiter(reject bool) RateLimiter {
	// TODO: replace with your actual RateLimiter
	return NewNoOpRateLimiter(reject)
}

func ProvideCommand() *cobra.Command {
	var (
		reject   bool
		msgCount int
		rate     int
		random   int
	)
	cmdRatelimit := &cobra.Command{
		Use:   "ratelimit",
		Short: "Run Rate Limiter",
		Long:  "Run ratelimiter and test output",
		Run: func(cmd *cobra.Command, args []string) {
			rl := ProvideRateLimiter(reject)
			err := rl.Start()
			if err != nil {
				log.Err(err).Msg("Failed to start RateLimiter")
				return
			}
			c := NewClient(rl, msgCount, rate, time.Millisecond*time.Duration(random))
			c.Run()
		},
	}
	cmdRatelimit.Flags().IntVarP(&msgCount, "messages", "m", 1000, "Number of messages to send")
	cmdRatelimit.Flags().IntVarP(&rate, "rate", "r", 100, "Messages sent per minute")
	cmdRatelimit.Flags().IntVarP(&random, "random", "a", 25, "Max amount of randomness (ms) added to rate")
	cmdRatelimit.Flags().BoolVarP(&reject, "reject", "j", false, "Reject messages rather than delay")

	return cmdRatelimit
}
