// Package secondly is a configuration management plugin for Go projects.
// It takes care of the app's configuration, specifically of updating it in
// runtime. It is capable of listening to file system events and signals to
// trigger configuration reload. It also provides a nice web GUI editor.
package secondly

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"reflect"
	"strings"
	"syscall"

	"github.com/howeyc/fsnotify"
)

var (
	config      interface{} // config stores application config
	configFile  string
	callbacks   = make(map[string][]func(oldVal, newVal interface{}))
	initialized bool
	initFunc    func()
)

// SetupFlags sets up Secondly's configuration flags.
func SetupFlags() {
	if flag.Parsed() {
		log.Fatalln("secondly.SetupFlags() must be called before flag.Parse()")
	}

	flag.StringVar(&configFile, "config", "config.json", "Path to config file")
}

// Manage accepts a pointer to a configuration struct.
func Manage(target interface{}) {
	if ok := isStructPtr(target); !ok {
		panic("Argument must be a pointer to a struct")
	}

	assign(target)

	bootstrap()
}

// StartServer will start an HTTP server with web interface to edit config.
func StartServer(host string, port int) {
	go startServer(fmt.Sprintf("%s:%d", host, port))
}

// HandleSIGHUP waits a SIGHUP system call and reloads configuration when
// receives one.
func HandleSIGHUP() {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGHUP)
	go func() {
		for range ch {
			log.Println("SIGHUP received, reloading config")
			readConfig()
		}
	}()
}

// HandleFileSystemEvents listens to file system events and reloads configuration when
// config file is modified.
func HandleFileSystemEvents() {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		panic(err)
	}
	if err := watcher.WatchFlags(filepath.Dir(configFile), fsnotify.FSN_MODIFY); err != nil {
		panic(err)
	}

	fname := configFile
	if ss := strings.Split(configFile, "/"); len(ss) > 1 {
		fname = ss[len(ss)-1]
	}

	go func() {
		for {
			select {
			case e := <-watcher.Event:
				if e.Name != fname {
					continue
				}
				log.Println("Config file was modified, reloading")
				readConfig()
			case err := <-watcher.Error:
				log.Println("fsnotify error:", err)
			}
		}
	}()
}

// OnLoad sets up a callback function that would be called once configuration
// is loaded for the first time.
func OnLoad(fun func()) {
	initFunc = fun
}

// OnChange adds a callback function that is triggered every time a value of
// a field changes. Field must be a json tag of the struct field.
func OnChange(field string, fun func(oldVal, newVal interface{})) {
	callbacks[field] = append(callbacks[field], fun)
}

// asign is responsible for assigning new config value. It is complicated
// because we're changing the value of an interface which is defined in
// another package.
func assign(target interface{}) {
	if config == nil {
		config = target
		return
	}

	cval := reflect.ValueOf(config).Elem()
	tval := reflect.ValueOf(target).Elem()
	cval.Set(tval)
}

// bootstrap sets up initial configuration.
func bootstrap() {
	if configFile == "" {
		log.Fatalln("path to config file is not set")
	}
	if fileExist(configFile) {
		log.Println("Loading config file")
		readConfig()
	} else {
		log.Fatalln("Config file not found")
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
	if err := writeFile(configFile, marshal(config)); err != nil {
		panic(err)
	}
}

func updateConfig(body []byte) {
	// Making a copy of old config for further comparison
	old := duplicate(config)
	// Making a second copy that we will fill with new data
	dupe := duplicate(config)

	if err := json.Unmarshal(body, dupe); err != nil {
		panic("Failed to update config")
		return
	}

	// Setting new config
	assign(dupe)

	triggerCallbacks(old, dupe)
}

func marshal(obj interface{}) []byte {
	body, err := json.Marshal(config)
	if err != nil {
		panic(err)
	}
	out := bytes.NewBuffer([]byte{})

	// Indent with empty prefix and four spaces
	if err = json.Indent(out, body, "", "    "); err != nil {
		panic(err)
	}

	// Adding a trailing newline
	// It's good for your carma
	out.WriteByte('\n')

	return out.Bytes()
}

func triggerCallbacks(oldConf, newConf interface{}) {
	// Don't trigger callbacks on fist load
	if !initialized {
		initialized = true
		if initFunc != nil {
			initFunc()
		}
		return
	}

	if len(callbacks) == 0 {
		return
	}

	for fname, d := range diff(oldConf, newConf) {
		if cbs, ok := callbacks[fname]; ok {
			for _, cb := range cbs {
				cb(d[0], d[1])
			}
		}
	}

	return
}

// Secondly accepts only a pointer to stuct as its config value. Here we're
// making sure the right argument is provided.
func isStructPtr(target interface{}) bool {
	if val := reflect.ValueOf(target); val.Kind() == reflect.Ptr {
		if val = reflect.Indirect(val); val.Kind() == reflect.Struct {
			return true
		}
	}

	return false
}

// duplicate creates a copy of a value behind the config interface. Such copies
// are used to replace config value and to check for changes.
func duplicate(original interface{}) interface{} {
	// Get the interface value
	val := reflect.ValueOf(original)
	// We expect a pointer to a struct, so now we need the underlying staruct
	val = reflect.Indirect(val)
	// Now we need the type (name) of this struct
	typ := val.Type()
	// Creating a duplicate instance of that struct
	dupe := reflect.New(typ)
	// Value copy
	dupe.Elem().Set(val)

	return dupe.Interface()
}
