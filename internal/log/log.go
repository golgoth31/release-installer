// Package log ...
package log

import (
	"github.com/rs/zerolog"
)

// Logger ...
var Logger zerolog.Logger

// SetLogger ...
func SetLogger(logger *zerolog.Logger) {
	Logger = logger.
		With().
		Timestamp().
		Logger()
}
