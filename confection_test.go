package confection2

import (
	"encoding/json"
	"testing"
)

type testConf struct {
	AppName  string           `json:"app_name"`
	Version  float32          `json:"version"`
	Database testDatabaseConf `json:"database"`
}
type testDatabaseConf struct {
	Adapter  string `json:"adapter"`
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
}

const (
	goodJSON = `{"app_name": "Confection", "version": 1}`
	badJSON  = `{"app_name": "noooo...`
)

func TestIsStructPtr(t *testing.T) {
	if ok := isStructPtr(1); ok {
		t.Error("Integer recognized as a struct pointer")
	}
	if ok := isStructPtr(testConf{}); ok {
		t.Error("Struct instance recognized as a struct pointer")
	}
	if ok := isStructPtr(&testConf{}); !ok {
		t.Error("Struct pointer was not recognized")
	}
}

func TestUnmarshal(t *testing.T) {
	conf := testConf{}
	var i interface{} = &conf

	if err := json.Unmarshal([]byte(badJSON), i); err == nil {
		t.Error("Expected error")
	}

	if err := json.Unmarshal([]byte(goodJSON), i); err != nil {
		t.Error("Unexpected error")
	}
	if conf.AppName != "Confection" {
		t.Errorf("Expected Foo to equal %q, got %q", "Confection", conf.AppName)
	}
	if conf.Version != 1 {
		t.Errorf("Expected Bar to equal %q, got %q", 1, conf.Version)
	}
}

func TestDuplicate(t *testing.T) {
	var i interface{} = &testConf{}

	dupe := duplicate(i)
	if _, ok := dupe.(*testConf); !ok {
		t.Error("Duplication failed")
	}
}
