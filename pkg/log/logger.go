package log

import (
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"reflect"
	"slices"
	"strings"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	globalInst *logger
	logLevel   zerolog.Level
)

type logger struct {
	writer io.Writer
}

func InitGlobalLogger(cfg *Config) {
	if globalInst == nil {
		writers := []io.Writer{}

		if slices.Contains(cfg.Targets, "file") {
			// File writer.
			fileWriter := &lumberjack.Logger{
				Filename:   cfg.Filename,
				MaxSize:    cfg.MaxSize,
				MaxBackups: cfg.MaxBackups,
				Compress:   cfg.Compress,
			}
			writers = append(writers, fileWriter)
		}

		if slices.Contains(cfg.Targets, "console") {
			// Console writer.
			writers = append(writers, zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: "15:04:05"})
		}

		globalInst = &logger{
			writer: io.MultiWriter(writers...),
		}

		// Set the global log level from the configuration.
		level, err := zerolog.ParseLevel(strings.ToLower(cfg.Level))
		if err != nil {
			level = zerolog.InfoLevel // Default to info level if parsing fails.
		}
		zerolog.SetGlobalLevel(level)

		log.Logger = zerolog.New(globalInst.writer).With().Timestamp().Logger()
	}
}

// NewLoggerLevel initializes the logger level.
func NewLoggerLevel(level zerolog.Level) {
	logLevel = level
}

// SetLoggerLevel sets logger level based on env.
func SetLoggerLevel(level string) {
	parsedLevel, err := zerolog.ParseLevel(strings.ToLower(level))
	if err != nil {
		parsedLevel = zerolog.InfoLevel // Default to info level if parsing fails
	}
	logLevel = parsedLevel
}

func GetCurrentLogLevel() zerolog.Level {
	return logLevel
}

func addFields(event *zerolog.Event, keyvals ...any) *zerolog.Event {
	if len(keyvals)%2 != 0 {
		keyvals = append(keyvals, "!MISSING-VALUE!")
	}

	for i := 0; i < len(keyvals); i += 2 {
		key, ok := keyvals[i].(string)
		if !ok {
			key = "!INVALID-KEY!"
		}

		value := keyvals[i+1]
		switch typ := value.(type) {
		case fmt.Stringer:
			if isNil(typ) {
				event.Any(key, typ)
			} else {
				event.Stringer(key, typ)
			}
		case error:
			event.AnErr(key, typ)
		case []byte:
			event.Str(key, hex.EncodeToString(typ))
		default:
			event.Any(key, typ)
		}
	}

	return event
}

func Trace(msg string, keyvals ...any) {
	addFields(log.Trace(), keyvals...).Msg(msg)
}

func Debug(msg string, keyvals ...any) {
	addFields(log.Debug(), keyvals...).Msg(msg)
}

func Info(msg string, keyvals ...any) {
	addFields(log.Info(), keyvals...).Msg(msg)
}

func Warn(msg string, keyvals ...any) {
	addFields(log.Warn(), keyvals...).Msg(msg)
}

func Error(msg string, keyvals ...any) {
	addFields(log.Error(), keyvals...).Msg(msg)
}

func Fatal(msg string, keyvals ...any) {
	addFields(log.Fatal(), keyvals...).Msg(msg)
}

func Panic(msg string, keyvals ...any) {
	addFields(log.Panic(), keyvals...).Msg(msg)
}

func isNil(i any) bool {
	if i == nil {
		return true
	}

	if reflect.TypeOf(i).Kind() == reflect.Ptr {
		return reflect.ValueOf(i).IsNil()
	}

	return false
}
