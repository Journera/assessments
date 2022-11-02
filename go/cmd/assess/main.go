package main

import (
	"github.com/journera/assessments/common"
	"github.com/journera/assessments/ratelimit"
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
)

func main() {
	log := common.ProvideLog()
	log.Info().Msg("Begin")

	var debug bool
	var rootCmd = &cobra.Command{
		Use:               "assess",
		CompletionOptions: cobra.CompletionOptions{DisableDefaultCmd: true},
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			if debug {
				zerolog.SetGlobalLevel(zerolog.DebugLevel)
			} else {
				zerolog.SetGlobalLevel(zerolog.InfoLevel)
			}
		},
	}
	rootCmd.PersistentFlags().BoolVarP(&debug, "debug", "d", false, "Enable debug logging")
	rootCmd.AddCommand(ratelimit.ProvideCommand())
	rootCmd.Execute()
}
