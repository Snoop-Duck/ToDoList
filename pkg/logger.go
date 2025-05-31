package logger

import (
	"os"
	"strconv"
	"sync"

	"github.com/rs/zerolog"
)

var log zerolog.Logger

var once sync.Once

func Get(debug ...bool) zerolog.Logger {
	once.Do(func() {
		zerolog.TimestampFieldName = "TIME"
		zerolog.LevelFieldName = "LEVEL"
		zerolog.CallerFieldName = "CALLER"
		zerolog.CallerMarshalFunc = func(_ uintptr, file string, line int) string {
			short := file
			for i := len(file) - 1; i > 0; i-- {
				if file[i] == '/' {
					short = file[i+1:]
					break
				}
			}
			return short + ":" + strconv.Itoa(line)
		}

		logLevel := zerolog.InfoLevel
		if len(debug) > 0 && debug[0] {
			logLevel = zerolog.DebugLevel
		}

		log = zerolog.New(os.Stdout).
			Level(logLevel).
			With().
			Timestamp().
			Caller().
			Logger().
			Output(zerolog.ConsoleWriter{Out: os.Stdout})

	})
	return log
}
