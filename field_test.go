package secondly

import (
	"testing"
)

func TestExtractFields(t *testing.T) {
	c := testConf{
		AppName: "Secondly",
		Version: 1.1,
		Database: testDatabaseConf{
			Adapter: "mysql",
			Host:    "localhost",
			Port:    3306,
		},
	}

	fields := indexFields(extractFields(c, ""))
	testField := func(fname, kind string, val interface{}) {
		if f, ok := fields[fname]; ok {
			if f.Kind != kind {
				t.Errorf("%s expected to be of kind %q, got %q", fname, kind, f.Kind)
			}
			if f.Value != val {
				t.Errorf("%s expected to have value %q, got %q", fname, val, f.Value)
			}
		} else {
			t.Errorf("Missing %s field", fname)
		}
	}

	testField("app_name", "string", c.AppName)
	testField("version", "float32", c.Version)
	testField("database.adapter", "string", c.Database.Adapter)
	testField("database.host", "string", c.Database.Host)
	testField("database.port", "int", c.Database.Port)
}

func TestDiff(t *testing.T) {
	c1 := testConf{
		AppName: "Secondly",
		Version: 1.3,
		Database: testDatabaseConf{
			Adapter: "mysql",
			Host:    "localhost",
			Port:    3306,
		},
	}
	c2 := testConf{
		AppName: "Secondly",
		Version: 2,
		Database: testDatabaseConf{
			Adapter:  "postgresql",
			Host:     "localhost",
			Port:     5432,
			Username: "root",
		},
	}

	d := diff(c1, c2)
	testField := func(fname string, oldVal, newVal interface{}) {
		if f, ok := d[fname]; ok {
			if f[0] != oldVal {
				t.Errorf("%s field old value was %q, not %q", oldVal, f[0])
			}
			if f[1] != newVal {
				t.Errorf("%s field new value was %q, not %q", newVal, f[1])
			}
		} else {
			t.Errorf("Expected %s field to have different values", fname)
		}
	}

	unchangedFields := []string{"app_name", "database.host", "database.password"}
	for _, f := range unchangedFields {
		if _, ok := d[f]; ok {
			t.Errorf("Expected %q field to be unchanged", f)
		}
	}

	testField("version", c1.Version, c2.Version)
	testField("database.adapter", c1.Database.Adapter, c2.Database.Adapter)
	testField("database.port", c1.Database.Port, c2.Database.Port)
	testField("database.username", c1.Database.Username, c2.Database.Username)
}
