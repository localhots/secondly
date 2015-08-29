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
