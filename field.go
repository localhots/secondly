package confection2

import (
	"log"
	"reflect"
)

type field struct {
	Kind string      `json:"kind"`
	Val  interface{} `json:"val"`
}

func extractFields(st interface{}, path string) map[string]field {
	res := make(map[string]field)
	typ := reflect.TypeOf(st)
	val := reflect.ValueOf(st)
	for i := 0; i < val.NumField(); i++ {
		ftyp := typ.Field(i)
		fval := val.Field(i)

		switch kind := fval.Kind(); kind {
		case reflect.Struct:
			sub := extractFields(fval.Interface(), ftyp.Name+".")
			for k, v := range sub {
				res[k] = v
			}
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
			res[path+ftyp.Name] = field{
				Kind: kind.String(),
				Val:  fval.Interface(),
			}
		default:
			log.Printf("Field type %q not supported for field %q\n", kind, path+ftyp.Name)
		}
	}

	return res
}
