package model

import "time"

// Event model variant.
type EventModel struct {
	Type string
}

// Make a event.
func MakeEvent(evType string, attributes map[string]any, timestamp time.Time) MeltModel {
	return MeltModel{
		Type:       Event,
		Timestamp:  timestamp.UnixMilli(),
		Attributes: attributes,
		Data: EventModel{
			Type: evType,
		},
	}
}
