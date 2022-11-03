package main

import (
	"github.com/journera/assessments/ratelimit"
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
)

func main() {
	var debug bool
	var trace bool
	var rootCmd = &cobra.Command{
		Use:               "assess",
		CompletionOptions: cobra.CompletionOptions{DisableDefaultCmd: true},
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			if trace {
				zerolog.SetGlobalLevel(zerolog.TraceLevel)
			} else if debug {
				zerolog.SetGlobalLevel(zerolog.DebugLevel)
			} else {
				zerolog.SetGlobalLevel(zerolog.InfoLevel)
			}
		},
	}
	rootCmd.PersistentFlags().BoolVarP(&debug, "debug", "d", false, "Enable debug logging")
	rootCmd.PersistentFlags().BoolVarP(&trace, "trace", "t", false, "Enable trace logging")
	rootCmd.AddCommand(ratelimit.ProvideRateLimitCommand())
	rootCmd.Execute()
}
