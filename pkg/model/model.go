package model

import (
	"encoding/json"
	"errors"
	"strconv"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/mitchellh/mapstructure"
)

type MeltType int

const (
	Metric MeltType = iota
	Event
	Log
	Trace
	Custom
)

// Intermediate model.
type MeltModel struct {
	Type MeltType
	// Unix timestamp in millis.
	Timestamp  int64
	Attributes map[string]any
	// Either a MetricModel, EventModel, LogModel, TraceModel or CustomModel.
	Data any
}

// Custom model variant.
type CustomModel struct {
	Id   string
	Data any
}

// Make a custom model.
func MakeCustom(id string, data any, timestamp time.Time) MeltModel {
	return MeltModel{
		Type:      Custom,
		Timestamp: timestamp.UnixMilli(),
		Data: CustomModel{
			Id:   id,
			Data: data,
		},
	}
}

func (receiver *MeltModel) UnmarshalJSON(data []byte) error {
	var dict map[string]any
	err := json.Unmarshal(data, &dict)
	if err != nil {
		return err
	}

	var model MeltModel
	err = mapstructure.Decode(dict, &model)
	if err != nil {
		return err
	}

	meltData := model.Data.(map[string]any)
	switch model.Type {
	case Metric:
		var metricModel MetricModel
		err := mapstructure.Decode(meltData, &metricModel)
		if err != nil {
			return err
		}
		model.Data = metricModel
	case Event:
		var eventModel EventModel
		err := mapstructure.Decode(meltData, &eventModel)
		if err != nil {
			return err
		}
		model.Data = eventModel
	case Log:
		var logModel LogModel
		err := mapstructure.Decode(meltData, &logModel)
		if err != nil {
			return err
		}
		model.Data = logModel
	case Trace:
		//TODO: unmarshal Trace model

	//TODO: Unmarshal Custom model

	default:
		return errors.New("'Type' contains an invalid value " + strconv.Itoa(int(model.Type)))
	}

	*receiver = model
	return nil
}

func (m *MeltModel) Metric() (MetricModel, bool) {
	model, ok := m.Data.(MetricModel)
	return model, ok
}

// Event obtains an EventModel from the MeltModel.
// If the inner data is a LogModel, it will be converted into an EventModel.
// This transformation may cause data loss: if the Log had a key in `Attributes` named "message", it will be overwritten
// with the contents of the `Message` field.
func (m *MeltModel) Event() (EventModel, bool) {
	if m.Type == Log {
		// Convert Log into an Event
		logModel := m.Data.(LogModel)
		model := EventModel{
			Type: logModel.Type,
		}
		if m.Attributes == nil {
			m.Attributes = map[string]any{}
		}
		// Warning: if the log already had an attribute named "message", it will be overwritten
		if _, ok := m.Attributes["message"]; ok {
			log.Warn("Log2Event: Log already had an attribute named 'message', overwriting it")
		}
		m.Attributes["message"] = logModel.Message
		return model, true
	} else {
		model, ok := m.Data.(EventModel)
		return model, ok
	}
}

// Log obtains a LogModel from the MeltModel.
// If the inner data is an EventModel, it will be converted into a LogModel.
// If the Event doesn't have a key in `Attributes` named "message", the `Message` field will remain empty.
func (m *MeltModel) Log() (LogModel, bool) {
	if m.Type == Event {
		message, ok := m.Attributes["message"].(string)
		if ok {
			delete(m.Attributes, "message")
		}
		eventModel, _ := m.Data.(EventModel)
		model := LogModel{
			Type:    eventModel.Type,
			Message: message,
		}
		return model, true
	} else {
		model, ok := m.Data.(LogModel)
		return model, ok
	}
}

func (m *MeltModel) Trace() (TraceModel, bool) {
	model, ok := m.Data.(TraceModel)
	return model, ok
}

func (m *MeltModel) Custom() (CustomModel, bool) {
	model, ok := m.Data.(CustomModel)
	return model, ok
}

// Numeric model.
type Numeric struct {
	IntOrFlt bool // true = Int, false = Float
	IntVal   int64
	FltVal   float64
}

// Numeric holds an integer.
func (n *Numeric) IsInt() bool {
	return n.IntOrFlt
}

// Numeric holds a float.
func (n *Numeric) IsFloat() bool {
	return !n.IntOrFlt
}

// Get float from Numeric.
func (n *Numeric) Float() float64 {
	if n.IsFloat() {
		return n.FltVal
	} else {
		return float64(n.IntVal)
	}
}

// Get int from Numeric.
func (n *Numeric) Int() int64 {
	if n.IsInt() {
		return n.IntVal
	} else {
		return int64(n.FltVal)
	}
}

// Get whatever it is.
func (n *Numeric) Value() any {
	if n.IsInt() {
		return n.IntVal
	} else {
		return n.FltVal
	}
}

// Make a Numeric from an int64.
func MakeIntNumeric(val int64) Numeric {
	return Numeric{
		IntOrFlt: true,
		IntVal:   val,
		FltVal:   0.0,
	}
}

// Make a Numeric from a float64.
func MakeFloatNumeric(val float64) Numeric {
	return Numeric{
		IntOrFlt: false,
		IntVal:   0,
		FltVal:   val,
	}
}
