package main

import (
	"bytes"
	"encoding/json"
	"github.com/buger/jsonparser"
	"github.com/pkg/errors"
	"log"
)

type Macro interface {
	ProcessFragment() ([]byte, error)
}

type macro struct {
	fragment  json.RawMessage
	functions map[string][]FunctionEvent
}

func NewMacro(fragment json.RawMessage) Macro {
	return &macro{
		fragment:  fragment,
		functions: make(map[string][]FunctionEvent),
	}
}

func (m *macro) ProcessFragment() ([]byte, error) {
	err := m.findFunctions()
	if err != nil {
		return nil, err
	}

	if len(m.functions) == 0 {
		log.Println("Template doesn't contain any Resource of type AWS::Serverless::Function that are using SharedHttpApi Event Type")
		return m.fragment, nil
	}
	for name, events := range m.functions {
		for _, event := range events {
			err := m.addPermissions(name, event)
			if err != nil {
				return nil, err
			}
			err = m.addIntegrationAndRoute(name, event)
			if err != nil {
				return nil, err
			}
		}
		m.deleteFunctionEvents(name)
	}
	log.Println(string(prettyPrint(m.fragment)))
	return m.fragment, nil
}

func (m *macro) addPermissions(funName string, event FunctionEvent) error {
	name, p, err := NewPermission(funName, event)
	if err == nil {
		m.fragment, err = jsonparser.Set(m.fragment, p, "Resources", name)
	}
	return err
}

func (m *macro) addIntegrationAndRoute(funName string, event FunctionEvent) error {
	integrationName, i, err := NewIntegration(funName, event)
	if err == nil {
		m.fragment, err = jsonparser.Set(m.fragment, i, "Resources", integrationName)
	}
	routeName, r, err := NewRoute(funName, integrationName, event)
	if err == nil {
		m.fragment, err = jsonparser.Set(m.fragment, r, "Resources", routeName)
	}
	return err
}

func (m *macro) deleteFunctionEvents(funName string) {
	m.fragment = jsonparser.Delete(m.fragment, "Resources", funName, "Properties", "Events")
}

// Find all AWS::Serverless::Function that are using SharedHttpApi Event type
func (m *macro) findFunctions() error {
	resources, _, _, err := jsonparser.Get(m.fragment, "Resources")

	err = jsonparser.ObjectEach(resources, func(key []byte, value []byte, dataType jsonparser.ValueType, offset int) error {
		rt, err := jsonparser.GetString(value, "Type")
		if err != nil {
			return err
		}
		if rt != "AWS::Serverless::Function" {
			return nil
		}

		events, err := getSharedHttpApiEvents(value)
		if err != nil {
			return err
		}
		if len(events) > 0 {
			m.functions[string(key)] = events
		}
		return nil
	})
	if err != nil {
		return errors.Wrap(err, "finding function resources")
	}
	return nil
}

// Check if AWS::Serverless::Function is having Events of Type: SharedHttpApi
func getSharedHttpApiEvents(function []byte) ([]FunctionEvent, error) {
	events, _, _, err := jsonparser.Get(function, "Properties", "Events")
	if err != nil {
		return nil, err
	}
	var functionEvents []FunctionEvent

	err = jsonparser.ObjectEach(events, func(key []byte, value []byte, dataType jsonparser.ValueType, offset int) error {
		et, err := jsonparser.GetString(value, "Type")
		if err != nil {
			return err
		}
		if et == "SharedHttpApi" {
			//json.Unmarshal()
			props, _, _, err := jsonparser.Get(value, "Properties")
			if err != nil {
				return err
			}
			var sharedHttpApi SharedHttpApi
			err = json.Unmarshal(props, &sharedHttpApi)
			if err != nil {
				return err
			}
			functionEvents = append(functionEvents, FunctionEvent{
				Name:          string(key),
				SharedHttpApi: sharedHttpApi,
			})
		}
		return nil
	})
	return functionEvents, nil
}

func prettyPrint(b []byte) []byte {
	var out bytes.Buffer
	err := json.Indent(&out, b, "", "  ")
	if err != nil {
		log.Println(err)
		return nil
	}
	return out.Bytes()
}
