package common

import (
	"github.com/rs/zerolog"
	"os"
)

var (
	log = createLogger()
)

func createLogger() *zerolog.Logger {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnixMicro
	output := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: "15:04:05.000"}
	l := zerolog.New(output).With().Timestamp().Logger()
	return &l
}

func ProvideLog() *zerolog.Logger {
	return log
}
