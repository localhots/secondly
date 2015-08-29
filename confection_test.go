package confection2

import (
	"encoding/json"
	"testing"
)

type testConf struct {
	Foo string `json:"foo"`
	Bar int    `json:"bar"`
}

const (
	goodJSON = `{"foo": "baz", "bar": 1}`
	badJSON  = `{"foo": "noooo...`
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
	if conf.Foo != "baz" {
		t.Errorf("Expected Foo to equal %q, got %q", "baz", conf.Foo)
	}
	if conf.Bar != 1 {
		t.Errorf("Expected Bar to equal %q, got %q", 1, conf.Bar)
	}
}

func TestDuplicate(t *testing.T) {
	var i interface{} = &testConf{}

	dupe := duplicate(i)
	if _, ok := dupe.(*testConf); !ok {
		t.Error("Duplication failed")
	}
}
