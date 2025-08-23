package logger

import (
	"bytes"
	"io"
	"os"
	"strconv"
	"sync"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

func TestGet_Initialization(t *testing.T) {
	t.Run("default info level", func(t *testing.T) {
		resetLogger()
		logger := Get()
		assert.Equal(t, zerolog.InfoLevel, logger.GetLevel())
	})

	t.Run("debug level when enabled", func(t *testing.T) {
		resetLogger()
		logger := Get(true)
		assert.Equal(t, zerolog.DebugLevel, logger.GetLevel())
	})

	t.Run("multiple calls return same instance", func(t *testing.T) {
		resetLogger()
		logger1 := Get()
		logger2 := Get(true)

		assert.Equal(t, logger1, logger2, "should return same instance")
		assert.Equal(t, zerolog.InfoLevel, logger2.GetLevel(), "level should not change after first init")
	})
}

func TestGet_OutputFormat(t *testing.T) {
	t.Run("console output format", func(t *testing.T) {
		resetLogger()
		old := os.Stdout
		r, w, _ := os.Pipe()
		os.Stdout = w

		logger := Get()
		logger.Info().Str("key", "value").Msg("test message")

		w.Close()
		os.Stdout = old

		var buf bytes.Buffer
		io.Copy(&buf, r)

		output := buf.String()

		// Проверяем наличие элементов, учитывая ANSI-коды
		assert.Contains(t, output, "test message")
		assert.Contains(t, output, "INF")
		assert.Contains(t, output, "key=")  // Проверяем часть строки до цветового кода
		assert.Contains(t, output, "value") // Проверяем значение
		assert.Contains(t, output, "logger_test.go")
	})
}

func TestCallerMarshalFunc(t *testing.T) {
	t.Run("windows path shortening", func(t *testing.T) {
		resetLogger()

		originalMarshalFunc := zerolog.CallerMarshalFunc

		zerolog.CallerMarshalFunc = func(_ uintptr, file string, line int) string {
			short := file
			for i := len(file) - 1; i > 0; i-- {
				if file[i] == '/' || file[i] == '\\' {
					short = file[i+1:]
					break
				}
			}
			return short + ":" + strconv.Itoa(line)
		}

		defer func() {
			zerolog.CallerMarshalFunc = originalMarshalFunc
		}()

		result := zerolog.CallerMarshalFunc(0, `C:\projects\app\file.go`, 42)
		assert.Equal(t, "file.go:42", result)
	})
}

func TestGet_Concurrency(t *testing.T) {
	resetLogger()
	var wg sync.WaitGroup
	iterations := 100

	results := make([]zerolog.Logger, iterations)

	for i := range iterations {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()
			results[index] = Get(index%2 == 0)
		}(i)
	}
	wg.Wait()

	first := results[0]
	for _, logger := range results {
		assert.Equal(t, first, logger)
	}
}

func TestCustomFieldNames(t *testing.T) {
	resetLogger()
	_ = Get()

	assert.Equal(t, "TIME", zerolog.TimestampFieldName)
	assert.Equal(t, "LEVEL", zerolog.LevelFieldName)
	assert.Equal(t, "CALLER", zerolog.CallerFieldName)
}

func resetLogger() {
	once = sync.Once{}
	log = zerolog.Logger{}
	zerolog.TimestampFieldName = "timestamp"
	zerolog.LevelFieldName = "level"
	zerolog.CallerFieldName = "caller"
	zerolog.CallerMarshalFunc = nil
}
