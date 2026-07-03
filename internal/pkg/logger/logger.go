package logger

import (
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

const bangkokTZ = "Asia/Bangkok"

func Init() {
	loc, err := time.LoadLocation(bangkokTZ)
	if err != nil {
		loc = time.UTC
	}

	zerolog.TimeFieldFormat = "02-Jan-2006 15:04:05.000"
	zerolog.TimestampFunc = func() time.Time {
		return time.Now().In(loc)
	}

	log.Logger = zerolog.New(os.Stdout).With().Timestamp().Caller().Logger()
}
