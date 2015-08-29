package main

import (
	"flag"
	"log"

	"github.com/localhots/secondly"
)

// testConf is our app's configuration
type testConf struct {
	AppName  string           `json:"app_name"`
	Version  float32          `json:"version"`
	Debug    bool             `json:"debug"`
	Database testDatabaseConf `json:"database"`
}

type testDatabaseConf struct {
	Adapter  string `json:"adapter"`
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
}

// conf is the variable that holds configuration
var conf = testConf{}

func main() {
	// Setting up flags
	secondly.SetupFlags()
	flag.Parse()

	// Delegating configuration management to Secondly
	secondly.Manage(&conf)
	// Handling file system events
	secondly.HandleFSEvents()
	// Handle SIGHUP
	secondly.HandleSIGHUP()
	// Starting a web server
	secondly.StartServer("", 5500)
	// Defining callbacks
	secondly.OnChange("app_name", func(o, n interface{}) {
		log.Printf("OMG! AppName changed from %q to %q", o, n)
	})

	// Other application startup logic
	select {}
}
