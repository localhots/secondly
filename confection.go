package confection2

import (
	"encoding/json"
	"flag"
	"log"
	"reflect"
)

var (
	config      interface{} // config stores application config
	configFile  string
	callbacks   = make(map[string][]func(oldVal, newVal interface{}))
	initialized bool
)

// SetFlags sets up Confection configuration flags.
func SetFlags() {
	flag.StringVar(&configFile, "config", "config.json", "Path to config file")
}

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

func bootstrap() {
	if configFile == "" {
		panic("path to config file is not set")
	}
	if fileExist(configFile) {
		log.Println("Loading config file")
		readConfig()
	} else {
		log.Println("Config file not found, saving an empty one")
		writeConfig()
	}
}

func readConfig() {
	body, err := readFile(configFile)
	if err != nil {
		panic(err)
	}
	updateConfig(body)
}

func writeConfig() {
	body, err := json.Marshal(config)
	if err != nil {
		panic(err)
	}
	if err = writeFile(configFile, body); err != nil {
		panic(err)
	}
}

func updateConfig(body []byte) {
	dupe := duplicate(config)
	if err := json.Unmarshal(body, dupe); err != nil {
		log.Println("Failed to update config")
		return
	}

	defer triggerCallbacks(config, dupe)

	// Setting new config
	config = dupe
}

func triggerCallbacks(oldConf, newConf interface{}) {
	// Don't trigger callbacks on fist load
	if !initialized {
		initialized = true
		return
	}

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
