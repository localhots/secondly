package confection2

import (
	"encoding/json"
	"log"
	"reflect"
)

var (
	// config stores application config.
	config    interface{}
	callbacks = make(map[string][]func(oldVal, newVal interface{}))
)

// Manage accepts a pointer to a configuration struct.
func Manage(target interface{}) {
	if ok := isStructPtr(target); !ok {
		panic("Argument must be a pointer to a struct")
	}

	config = target
}

// OnChange adds a callback function that is triggered every time a value of
// a field changes.
func OnChange(field string, fun func(oldVal, newVal interface{})) {
	callbacks[field] = append(callbacks[field], fun)
}

func update(body []byte) {
	dupe := duplicate(config)
	if err := unmarshal(body, dupe); err != nil {
		log.Println("Failed to update config")
		return
	}

	defer triggerCallbacks(config, dupe)

	// Setting new config
	config = dupe
}

func triggerCallbacks(oldConf, newConf interface{}) {
	return
}

func isStructPtr(target interface{}) bool {
	if val := reflect.ValueOf(target); val.Kind() == reflect.Ptr {
		if val = reflect.Indirect(val); val.Kind() == reflect.Struct {
			return true
		}
	}

	return false
}

func unmarshal(body []byte, target interface{}) error {
	if err := json.Unmarshal(body, target); err != nil {
		return err
	}

	return nil
}

func duplicate(original interface{}) interface{} {
	// Get the interface value
	val := reflect.ValueOf(original)
	// We expect a pointer to a struct, so now we need the underlying staruct
	val = reflect.Indirect(val)
	// Now we need the type (name) of this struct
	typ := val.Type()
	// Creating a duplicate instance of that struct
	dupe := reflect.New(typ).Interface()

	return dupe
}
