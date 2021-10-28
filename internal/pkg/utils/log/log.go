package log

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"io"
	"os"
	"time"
)

func init() {
	zerolog.SetGlobalLevel(zerolog.DebugLevel)
	SetOutput(os.Stdout)
}

func SetOutput(out io.Writer) {
	log.Output(zerolog.ConsoleWriter{
		Out:        out,
		TimeFormat: time.RFC3339,
		NoColor:    !(out == os.Stdout || out == os.Stderr),
	})
}

func Debug(msg string) {
	log.Debug().Msg(msg)
}

func Info(msg string) {
	log.Info().Msg(msg)
}

func Warn(msg string) {
	log.Warn().Msg(msg)
}

func Error(err error) {
	log.Error().Err(err).Msg("")
}
