package export

import (
	"fmt"
	"newrelic/multienv/pkg/config"
	"newrelic/multienv/pkg/model"
	"strconv"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/newrelic/infra-integrations-sdk/v4/data/event"
	"github.com/newrelic/infra-integrations-sdk/v4/data/metric"
	"github.com/newrelic/infra-integrations-sdk/v4/integration"
)

const NrInfraInventory = "NrInfraInventory"

// Inventory model to insert into a MeltModel (as Custom MeltType)
type NrInfraInventoryData struct {
	key   string
	field string
	value any
}

// Build a MeltModel with a custom NrInfraInventoryData
func MakeInventory(key string, field string, value any) model.MeltModel {
	data := NrInfraInventoryData{
		key:   key,
		field: field,
		value: value,
	}
	return model.MakeCustom(NrInfraInventory, data, time.Now())
}

func exportNrInfra(pipeConf config.PipelineConfig, data []model.MeltModel) error {
	log.Print("------> NR Infra Exporter = ", data)

	name, ok := pipeConf.GetString("name")
	if !ok {
		name = "InfraIntegration"
	}

	version, ok := pipeConf.GetString("version")
	if !ok {
		version = "0.1.0"
	}

	// Create integration
	i, err := integration.New(name, version)
	if err != nil {
		log.Error("Error creating Nr Infra integration", err)
		return err
	}

	entityName, ok := pipeConf.GetString("entity_name")
	if !ok {
		entityName = "EntityName"
	}

	entityType, ok := pipeConf.GetString("entity_type")
	if !ok {
		entityType = "EntityType"
	}

	entityDisplay, ok := pipeConf.GetString("entity_display")
	if !ok {
		entityDisplay = "EntityDisplay"
	}

	// Create entity
	entity, err := i.NewEntity(entityName, entityType, entityDisplay)
	if err != nil {
		log.Error("Error creating entity", err)
		return err
	}

	for _, d := range data {
		if d.Type == model.Event || d.Type == model.Log {
			ev, ok := d.Event()
			if ok {
				nriEv, err := event.New(time.UnixMilli(d.Timestamp), "Event of type "+ev.Type, ev.Type)
				nriEv.Attributes = d.Attributes
				if err != nil {
					log.Error("Error creating event", err)
				} else {
					entity.AddEvent(nriEv)
				}
			}
		} else if d.Type == model.Metric {
			m, ok := d.Metric()
			if ok {
				switch m.Type {
				case model.Gauge:
					gauge, err := integration.Gauge(time.UnixMilli(d.Timestamp), m.Name, m.Value.Float())
					addAttributes(&d, &gauge)
					if err != nil {
						log.Error("Error creating gauge metric", err)
					} else {
						entity.AddMetric(gauge)
					}
				case model.Count, model.CumulativeCount:
					//TODO: NO TIME INTERVAL???
					count, err := integration.Count(time.UnixMilli(d.Timestamp), m.Name, m.Value.Float())
					addAttributes(&d, &count)
					if err != nil {
						log.Error("Error creating count metric", err)
					} else {
						entity.AddMetric(count)
					}
				case model.Summary:
					//TODO
				}
			}
		} else if d.Type == model.Custom {
			inv, ok := d.Custom()
			if ok {
				if inv.Id == NrInfraInventory {
					dat, ok := inv.Data.(NrInfraInventoryData)
					if ok {
						entity.AddInventoryItem(dat.key, dat.field, dat.value)
					} else {
						log.Error("Custom data should be of type NrInfraInventoryData")
					}
				} else {
					log.Warn("Ignored data, not NrInfraInventory")
				}
			}
		} else {
			log.Warn("Ignored data, not a metric, event or log: ", d)
		}
	}

	i.AddEntity(entity)

	err = i.Publish()
	if err != nil {
		log.Error("Error publishing", err)
		return err
	}

	return nil
}

func addAttributes(model *model.MeltModel, metric *metric.Metric) {
	for k, v := range model.Attributes {
		switch val := v.(type) {
		case string:
			(*metric).AddDimension(k, val)
		case int:
			(*metric).AddDimension(k, strconv.Itoa(val))
		case float32:
			(*metric).AddDimension(k, strconv.FormatFloat(float64(val), 'f', 2, 32))
		case float64:
			(*metric).AddDimension(k, strconv.FormatFloat(val, 'f', 2, 32))
		case fmt.Stringer:
			(*metric).AddDimension(k, val.String())
		default:
			log.Warn("Attribute of unsupported type: ", k, v)
		}

	}
}
