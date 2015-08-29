package confection2

import (
	"encoding/json"
	"reflect"
)

var (
	// config stores application config.
	config interface{}
)

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
