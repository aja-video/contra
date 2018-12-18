package collectors

import (
	"log"
	"os"
	"reflect"
)

type CollectorFactory struct {
	registry map[string]reflect.Type
}

///
/// Everything was working pretty well, except got incredibly stuck
/// trying to take the reflection created struct to the interface.
/// DeviceComware can't be an interface because you can't set functions
/// on interfaces. But, if it's a struct, although reflect.New appears
/// to create the struct, and the struct appears to have the functions.
/// Trying to cast the struct.(Collector) insists that the functions are
/// not declared.
///
/// Hopefully I'm just missing something that I'll see tomorrow.
///

func (cf *CollectorFactory) Register(name string, i interface{}) {
	// Ensure map is initialized
	if cf.registry == nil {
		cf.registry = make(map[string]reflect.Type)
	}
	// Lookup the type.
	t := reflect.TypeOf(i)
	cf.registry[t.Name()] = t
	cf.registry[name] = t
	// Debugging
	log.Printf("Registered type: [ %s ] With alias: [ %s ]", t.Name(), name)

	n := Collector(&DeviceVyatta{})
	log.Println(n)
	nn := n.(Collector)
	log.Println(nn)

	n = Collector(&DeviceProcurve{})
	log.Println(n)
	nn = n.(Collector)
	log.Println(nn)

	//n = Collector(&DeviceCiscoCsb{})
	//o := Collector(new(DeviceCiscoCsb{}))
	//log.Println(o)
	//nn = o.(Collector)
	//log.Println(nn)

	//v := reflect.New(cf.registry[name]).Elem()
	vat := reflect.New(cf.registry[name])
	mat := reflect.Indirect(vat).Interface().(Collector)
	//y := v.Interface()
	//z := Collector(new(y))
	//z := y.(Collector)
	log.Println(mat)
	os.Exit(0)
}

func (cf *CollectorFactory) Make(typeName string) interface{} {
	log.Println(cf.registry)
	log.Println(typeName)
	log.Println(cf.registry[typeName])
	// return nil, fmt.Errorf("unrecognized collector type: %v", d.Type)

	v := reflect.New(cf.registry[typeName]).Elem()
	y := v.Interface()
	z := y.(Collector)
	log.Println(z)

	return v.Interface()
}

// MakeCollector will generate the appropriate collector based on the
// type string passed in by the configuration.
func (cf *CollectorFactory) MakeCollector(typeName string) Collector {
	// Same as Make, but casted to collector for convenience.
	// Would be nice if we could assign the DeviceConfig here, but it
	// isn't clear how to do that in golang.
	return cf.Make(typeName).(Collector)
}
