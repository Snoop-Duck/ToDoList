package logger

import (
	"fmt"
	"os"
	"strconv"
	"sync"

	"github.com/rs/zerolog"
)

var log zerolog.Logger

var once sync.Once

func Get(flag ...bool) zerolog.Logger {
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
			file = short
			return file + ":" + strconv.Itoa(line)
		}
		if flag[0] {
			log = zerolog.New(os.Stdout).Level(zerolog.DebugLevel).With().Timestamp().Caller().Logger().Output(zerolog.ConsoleWriter{Out: os.Stdout})
		} else {
			log = zerolog.New(os.Stdout).Level(zerolog.InfoLevel).With().Timestamp().Caller().Logger().Output(zerolog.ConsoleWriter{Out: os.Stdout})
		}

	})
	fmt.Println("return logger")
	return log
}
