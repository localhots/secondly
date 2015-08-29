package confection2

import (
	"testing"
)

func TestExtractFields(t *testing.T) {
	c := testConf{
		AppName: "Confection",
		Version: 1.1,
		Database: testDatabaseConf{
			Adapter: "mysql",
			Host:    "localhost",
			Port:    3306,
		},
	}

	fields := extractFields(c, "")
	testField := func(fieldName, kind string, val interface{}) {
		if f, ok := fields[fieldName]; ok {
			if f.Kind != kind {
				t.Errorf("%s expected to be of kind %q, got %q", fieldName, kind, f.Kind)
			}
			if f.Val != val {
				t.Errorf("%s expected to have value %q, got %q", fieldName, val, f.Val)
			}
		} else {
			t.Errorf("Missing %s field", fieldName)
		}
	}

	testField("AppName", "string", c.AppName)
	testField("Version", "float32", c.Version)
	testField("Database.Adapter", "string", c.Database.Adapter)
	testField("Database.Host", "string", c.Database.Host)
	testField("Database.Port", "int", c.Database.Port)
}
