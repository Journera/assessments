package ratelimit

import (
	"github.com/brianvoe/gofakeit"
	"github.com/spf13/cobra"
	"sync"
	"time"
)

var (
	reject      bool
	clientCount int
	msgCount    int
	sendRate    int
	limitRate   int
	variance    int
)

func ProvideRateLimitCommand() *cobra.Command {
	cmdRatelimit := &cobra.Command{
		Use:   "ratelimit",
		Short: "Run Rate Limiter",
		Long:  "Run Rate Limiter and test output",
		Run: func(cmd *cobra.Command, args []string) {
			start := time.Now()
			log.Info().Int("Limit", limitRate).Msg("RateLimit Begin")
			limiter := ProvideRateLimiter(limitRate, reject)
			err := limiter.Start()
			if err != nil {
				log.Err(err).Msg("Failed to start RateLimiter")
				return
			}

			var done sync.WaitGroup
			done.Add(clientCount)
			go func() {
				done.Wait()     // wait for all senders to complete
				limiter.Close() // stop the limited to end the collector
			}()

			gofakeit.Seed(time.Now().Unix())
			collector := NewCollector(limiter)

			log.Info().Msgf("Starting %d clients", clientCount)
			senders := createSenders(clientCount)
			for i, sender := range senders {
				c := NewClient(limiter, sender, msgCount, calculateRate(clientCount, i, sendRate, variance))
				go func() { // run each sender
					c.Run()
					done.Done()
				}()
			}
			collector.Run()
			Evaluate(collector)
			log.Info().Stringer("RunTime", time.Since(start)).Msg("RateLimit Complete |")
		},
	}
	cmdRatelimit.Flags().IntVarP(&clientCount, "clients", "c", 5, "Number of clients")
	cmdRatelimit.Flags().IntVarP(&limitRate, "limit-rate", "l", 100, "Rate limit of messages per minute per client")
	cmdRatelimit.Flags().IntVarP(&sendRate, "send-rate", "s", 100, "Messages sent per minute per client")
	cmdRatelimit.Flags().IntVarP(&variance, "variance", "v", 10, "Msg rate variance per client")
	cmdRatelimit.Flags().IntVarP(&msgCount, "messages", "m", 1000, "Number of messages to send")
	cmdRatelimit.Flags().BoolVarP(&reject, "reject", "r", false, "Reject messages rather than delay")

	return cmdRatelimit
}

// createSenders creates a unique list of random names as the senders
func createSenders(count int) []string {
	senders := make([]string, 0, count)
	for i := 0; i < count; i++ {
		for {
			found := false
			name := gofakeit.FirstName()
			for j := 0; j < len(senders); j++ {
				if senders[j] == name {
					found = true
				}
			}
			if !found {
				senders = append(senders, name)
				break
			}
		}
	}
	return senders
}

// calculateRate will give a varied send rate per client. This will attempt to balance
// the values on either side of the target rate.
// Example: clients: 5, rate: 100, variance: 10
// Results: 80, 90, 100, 110, 120
func calculateRate(numClients, seq, rate, variance int) int {
	if numClients == 1 {
		return rate + variance
	}
	i := numClients / 2 // integer math is desired here
	return ((seq - i) * variance) + rate
}

func Evaluate(collector *Collector) {
	log.Info().Msgf("Evaluate | %d messages received", collector.received.Size())

	bySender := make(map[string]*Stats)
	for msg := range collector.received.Iter() {
		stats, ok := bySender[msg.Sender]
		if ok {
			stats.Messages.AddLast(msg)
		} else {
			stats = NewStats(msg)
			bySender[msg.Sender] = stats
		}
	}

	if len(bySender) != clientCount {
		log.Error().Msgf("The number of clients (%d) does not match with the received messages (%d)", clientCount, len(bySender))
	}

	totalMsgs := 0
	for sender, stats := range bySender {
		log.Info().Msgf("== Checking sender %s ==", sender)
		totalMsgs += stats.Messages.Size()
		for msg := range stats.Messages.Iter() {
			dur := msg.ReceiveTime.Sub(msg.SendTime)
			stats.TotalTime += dur
			if stats.MinTime == 0 || dur < stats.MinTime {
				stats.MinTime = dur
			}
			if dur > stats.MaxTime {
				stats.MaxTime = dur
			}
		}
		stats.AverageTime = stats.TotalTime / time.Duration(stats.Messages.Size())
		log.Info().
			Int("Msgs", stats.Messages.Size()).
			Dur("Min", stats.MinTime).
			Dur("Max", stats.MaxTime).
			Dur("Avg", stats.AverageTime).
			Msg("Stats |")

		if stats.Messages.Size() < msgCount {
			log.Warn().Int("Msgs", stats.Messages.Size()).
				Msgf("%d messages were lost", msgCount-stats.Messages.Size())
		}
	}
	log.Info().Msg("===========")
	log.Info().
		Int("Sent", msgCount*clientCount).
		Int("Rcvd", totalMsgs).
		Msg("Messages |")
}
