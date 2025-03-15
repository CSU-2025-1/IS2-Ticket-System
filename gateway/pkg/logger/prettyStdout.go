package logger

import (
	"github.com/rs/zerolog"
	"os"
	"time"
)

const (
	Debug = iota
	Info
	Warning
	Error
	Trace
)

// PrettyStdout logging pretty messages into os.Stdout
type PrettyStdout struct {
	logger zerolog.Logger
}

// NewPrettyStdout returns new StdOutLogger
func NewPrettyStdout(level int) *PrettyStdout {
	return &PrettyStdout{
		logger: zerolog.New(
			zerolog.ConsoleWriter{
				Out:        os.Stdout,
				TimeFormat: time.RFC3339,
			},
		).Level(zerolog.Level(level)),
	}
}

// Infof logs INFO level messages
func (p *PrettyStdout) Infof(message string, args ...interface{}) {
	p.logger.Info().Timestamp().Msgf(message, args...)
}

// Warnf logs WARN level messages
func (p *PrettyStdout) Warnf(msg string, args ...interface{}) {
	p.logger.Warn().Timestamp().Msgf(msg, args...)
}

// Errorf logs ERROR level messages
func (p *PrettyStdout) Errorf(msg string, args ...interface{}) {
	p.logger.Error().Timestamp().Msgf(msg, args...)
}
