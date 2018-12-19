package collectors

import (
	"fmt"
	"reflect"
)

type collectorFactory struct {
	registry map[string]reflect.Type
}

func (cf *collectorFactory) Register(name string, i interface{}) {
	// Ensure map is initialized
	if cf.registry == nil {
		cf.registry = make(map[string]reflect.Type)
	}

	// Lookup the type.
	t := reflect.TypeOf(i)
	cf.registry[t.Name()] = t
	cf.registry[name] = t

	// Debugging
	//log.Printf("Registered type: [ %s ] With alias: [ %s ]", t.Name(), name)
}

func (cf *collectorFactory) Make(typeName string) (interface{}, error) {
	if val, ok := cf.registry[typeName]; ok {
		return reflect.New(val).Elem().Addr().Interface(), nil
	}

	return nil, fmt.Errorf("unrecognized collector type: %v", typeName)
}

// MakeCollector will generate the appropriate collector based on the
// type string passed in by the configuration.
func (cf *collectorFactory) MakeCollector(typeName string) (Collector, error) {
	// Same as Make, but casted to collector for convenience.
	// Would be nice if we could assign the DeviceConfig here, but it
	// isn't clear how to do that in golang.
	c, err := cf.Make(typeName)
	if err != nil {
		return nil, err
	}

	return c.(Collector), nil
}
