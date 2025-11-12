package log

import (
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func SetUp(level string, nocolor bool) {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnixMs
	console := zerolog.ConsoleWriter{Out: os.Stderr, NoColor: nocolor}
	log.Logger = zerolog.New(console).With().Timestamp().Logger()
	lvl, err := zerolog.ParseLevel(level)
	if err != nil {
		lvl = zerolog.InfoLevel
	}
	zerolog.SetGlobalLevel(lvl)
}
