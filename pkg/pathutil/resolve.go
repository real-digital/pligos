package pathutil

import (
	"path/filepath"
	"reflect"
)

func Resolve(in interface{}, root string) {
	t := reflect.TypeOf(in).Elem()
	val := reflect.ValueOf(in).Elem()

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)

		switch val.Field(i).Kind() {
		case reflect.String:
			if field.Tag.Get("filepath") == "resolve" {
				val.Field(i).SetString(filepath.Join(root, val.Field(i).String()))
			}
		case reflect.Struct:
			next := val.Field(i).Addr()
			Resolve(next.Interface(), root)

			val.Field(i).Set(reflect.ValueOf(next.Elem().Interface()))
		case reflect.Slice:
			if val.Field(i).Len() == 0 {
				continue
			}

			switch val.Field(i).Index(0).Kind() {
			case reflect.Struct:
				for j := 0; j < val.Field(i).Len(); j++ {
					next := val.Field(i).Index(j).Addr()
					Resolve(next.Interface(), root)

					val.Field(i).Index(j).Set(reflect.ValueOf(next.Elem().Interface()))
				}
			case reflect.String:
				for j := 0; j < val.Field(i).Len(); j++ {
					if field.Tag.Get("filepath") == "resolve" {
						val.Field(i).Index(j).SetString(filepath.Join(root, val.Field(i).Index(j).String()))
					}
				}

			}

		}

	}
}
