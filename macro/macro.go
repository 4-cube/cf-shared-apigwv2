package macro

import (
	"encoding/json"
	"github.com/4-cube/cf-shared-apigwv2/pkg/cf"
	"github.com/buger/jsonparser"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type Macro interface {
	ProcessFragment() ([]byte, error)
}

type macro struct {
	fragment       json.RawMessage
	log            *logrus.Logger
	functionEvents map[string][]*cf.HttpApiEvent
}

func NewMacro(fragment json.RawMessage, log *logrus.Logger) Macro {
	return &macro{
		fragment:       fragment,
		log:            log,
		functionEvents: make(map[string][]*cf.HttpApiEvent),
	}
}

func (m *macro) ProcessFragment() ([]byte, error) {
	err := m.findFunctions()
	if err != nil {
		m.log.Errorln(err)
		return nil, err
	}

	if len(m.functionEvents) == 0 {
		m.log.Warn("Template doesn't contain any Resource of type AWS::Serverless::Function that is using HttpApi Event Type")
		return m.fragment, nil
	}

	for fnName, events := range m.functionEvents {
		m.deleteFunctionEvents(fnName)
		for _, event := range events {
			err := m.createPermission(fnName, event)
			if err != nil {
				return nil, err
			}
			err = m.createIntegrationAndRoute(fnName, event)
			if err != nil {
				return nil, err
			}
		}
	}
	return m.fragment, nil
}

func (m *macro) createPermission(fnName string, event *cf.HttpApiEvent) error {
	b := cf.NewPermissionBuilder(fnName, event)
	m.log.Infof("Creating Resource %s", b.Name())
	res, err := b.JSON()
	if err != nil {
		m.log.Errorln(err)
		return errors.Wrap(err, "creating permission resource failed")
	}
	m.fragment, err = jsonparser.Set(m.fragment, res, cf.ResourceKey, b.Name())
	if err != nil {
		m.log.Errorln(err)
		return errors.Wrap(err, "adding permission resource failed")
	}
	return nil
}

func (m *macro) createIntegrationAndRoute(fnName string, event *cf.HttpApiEvent) error {
	ib := cf.NewIntegrationBuilder(fnName, event)

	m.log.Infof("Creating Resource %s", ib.Name())

	ires, err := ib.JSON()
	if err != nil {
		m.log.Errorln(err)
		return errors.Wrap(err, "creating integration resource failed")
	}

	m.fragment, err = jsonparser.Set(m.fragment, ires, cf.ResourceKey, ib.Name())
	if err != nil {
		m.log.Errorln(err)
		return errors.Wrap(err, "adding integration resource failed")
	}

	rb := cf.NewRouteBuilder(fnName, ib.Name(), event)

	m.log.Infof("Creating Resource %s", rb.Name())

	rres, err := rb.JSON()
	if err != nil {
		m.log.Errorln(err)
		return errors.Wrap(err, "creating route resource failed")
	}

	m.fragment, err = jsonparser.Set(m.fragment, rres, cf.ResourceKey, rb.Name())
	if err != nil {
		m.log.Errorln(err)
		return errors.Wrap(err, "adding route resource failed")
	}

	return nil
}

func (m *macro) deleteFunctionEvents(fnName string) {
	m.fragment = jsonparser.Delete(m.fragment, cf.ResourceKey, fnName, cf.PropertiesKey, cf.EventsKey)
}

// Find all AWS::Serverless::Function that are using HttpApi Event type
func (m *macro) findFunctions() error {
	resources, _, _, err := jsonparser.Get(m.fragment, cf.ResourceKey)

	err = jsonparser.ObjectEach(resources, func(key []byte, value []byte, dataType jsonparser.ValueType, offset int) error {
		rt, err := jsonparser.GetString(value, cf.TypeKey)
		if err != nil {
			return err
		}
		if rt != cf.FunctionResourceType {
			return nil
		}

		events, err := findEvents(value)
		if err != nil {
			return err
		}
		if len(events) > 0 {
			m.addFunctionEvents(string(key), events)
		}
		return nil
	})
	if err != nil {
		return errors.Wrap(err, "finding function resources")
	}
	return nil
}

// Find all of AWS::Serverless::Function Events of Type: HttpApi
func findEvents(function []byte) ([]*cf.HttpApiEvent, error) {
	events, _, _, err := jsonparser.Get(function, cf.PropertiesKey, cf.EventsKey)
	if err != nil {
		return nil, err
	}

	var httpApiEvents []*cf.HttpApiEvent

	err = jsonparser.ObjectEach(events, func(key []byte, value []byte, dataType jsonparser.ValueType, offset int) error {
		et, err := jsonparser.GetString(value, cf.TypeKey)
		if err != nil {
			return err
		}
		if et == cf.HttpApiType {
			httpApiEvent, err := unmarshalEvent(string(key), value)
			if err != nil {
				return err
			}
			httpApiEvents = append(httpApiEvents, httpApiEvent)
		}
		return nil
	})
	return httpApiEvents, nil
}

func unmarshalEvent(key string, j []byte) (*cf.HttpApiEvent, error) {
	props, _, _, err := jsonparser.Get(j, cf.PropertiesKey)
	if err != nil {
		return nil, err
	}
	var httpApi cf.HttpApi
	err = json.Unmarshal(props, &httpApi)
	if err != nil {
		return nil, err
	}
	return &cf.HttpApiEvent{
		Name:    key,
		HttpApi: httpApi,
	}, nil
}

func (m *macro) addFunctionEvents(fnName string, events []*cf.HttpApiEvent) {
	m.functionEvents[fnName] = events
}
