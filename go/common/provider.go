package common

import (
	"fmt"
	"github.com/rs/zerolog"
	"os"
	"strings"
)

var (
	log = createLogger()
)

func createLogger() *zerolog.Logger {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnixMicro
	output := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: "15:04:05.000"}
	output.FormatLevel = func(i interface{}) string {
		return strings.ToUpper(fmt.Sprintf("| %-6s|", i))
	}
	output.FormatMessage = func(i interface{}) string {
		return fmt.Sprintf("%s", i)
	}
	output.FormatFieldName = func(i interface{}) string {
		return fmt.Sprintf("%s:", i)
	}
	output.FormatFieldValue = func(i interface{}) string {
		return fmt.Sprintf("%s", i)
	}
	l := zerolog.New(output).With().Timestamp().Logger()
	return &l
}

func ProvideLog() *zerolog.Logger {
	return log
}
