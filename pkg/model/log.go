package model

import (
	"time"
)

// Log model variant.
type LogModel struct {
	Message string
	Type    string
}

// Make a log.
func MakeLog(message string, logType string, timestamp time.Time) MeltModel {
	return MeltModel{
		Type:      Log,
		Timestamp: timestamp.UnixMilli(),
		Data: LogModel{
			Message: message,
			Type:    logType,
		},
	}
}
