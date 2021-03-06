package secondly

import (
	"log"
	"reflect"
)

type field struct {
	Path  string      `json:"path"`
	Name  string      `json:"name"`
	Kind  string      `json:"kind"`
	Value interface{} `json:"value"`
}

func extractFields(st interface{}, path string) []field {
	var res []field

	val := reflect.ValueOf(st)
	if val.Kind() == reflect.Ptr {
		val = reflect.Indirect(val)
	}
	typ := val.Type()

	for i := 0; i < val.NumField(); i++ {
		ftyp := typ.Field(i)
		fval := val.Field(i)
		tag := ftyp.Tag.Get("json")

		switch kind := fval.Kind(); kind {
		case reflect.Struct:
			sub := extractFields(fval.Interface(), path+tag+".")
			res = append(res, sub...)
		case reflect.Bool,
			reflect.Int,
			reflect.Int8,
			reflect.Int16,
			reflect.Int32,
			reflect.Int64,
			reflect.Uint,
			reflect.Uint8,
			reflect.Uint16,
			reflect.Uint32,
			reflect.Uint64,
			reflect.Float32,
			reflect.Float64,
			reflect.String:
			res = append(res, field{
				Path:  path + tag,
				Name:  ftyp.Name,
				Kind:  kind.String(),
				Value: fval.Interface(),
			})
		default:
			log.Printf("Field type %q not supported for field %q\n", kind, path+tag)
		}
	}

	return res
}

func diff(a, b interface{}) map[string][]interface{} {
	af := indexFields(extractFields(a, ""))
	bf := indexFields(extractFields(b, ""))

	res := make(map[string][]interface{})
	for name, f := range af {
		if bf[name].Value != f.Value {
			res[name] = []interface{}{f.Value, bf[name].Value}
		}
	}

	return res
}

func indexFields(fields []field) map[string]field {
	res := make(map[string]field)
	for _, f := range fields {
		res[f.Path] = f
	}

	return res
}
