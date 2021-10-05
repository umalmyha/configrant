package structs

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"
)

type FieldSetter interface {
	Apply(field reflect.Value, value string) error
}

func determineFieldSetter(typ reflect.Type) (setter FieldSetter, err error) {
	switch typ.Kind() {
	case reflect.String:
		setter = new(stringFieldSetter)
	case reflect.Bool:
		setter = new(boolFieldSetter)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if isTimeDurationType(typ) {
			setter = new(timeDurationFieldSetter)
		} else {
			setter = new(intFieldSetter)
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		setter = new(uintFieldSetter)
	case reflect.Float32, reflect.Float64:
		setter = new(floatFieldSetter)
	case reflect.Slice:
		setter = new(sliceFieldSetter)
	case reflect.Map:
		setter = new(mapFieldSetter)
	default:
		err = fmt.Errorf("type %s is not supported for configuration", typ.Name())
	}
	return
}

func isTimeDurationType(typ reflect.Type) bool {
	return typ.Kind() == reflect.Int64 && typ.PkgPath() == "time" && typ.Name() == "Duration"
}

type stringFieldSetter struct{}

func (s *stringFieldSetter) Apply(field reflect.Value, value string) error {
	field.SetString(value)
	return nil
}

type boolFieldSetter struct{}

func (s *boolFieldSetter) Apply(field reflect.Value, value string) error {
	if boolValue, err := strconv.ParseBool(value); err != nil {
		return err
	} else {
		field.SetBool(boolValue)
		return nil
	}
}

type intFieldSetter struct{}

func (s *intFieldSetter) Apply(field reflect.Value, value string) error {
	if intValue, err := strconv.ParseInt(value, 0, field.Type().Bits()); err != nil {
		return err
	} else {
		field.SetInt(intValue)
		return nil
	}
}

type uintFieldSetter struct{}

func (s *uintFieldSetter) Apply(field reflect.Value, value string) error {
	if uintValue, err := strconv.ParseUint(value, 0, field.Type().Bits()); err != nil {
		return err
	} else {
		field.SetUint(uintValue)
		return nil
	}
}

type floatFieldSetter struct{}

func (s *floatFieldSetter) Apply(field reflect.Value, value string) error {
	if floatValue, err := strconv.ParseFloat(value, field.Type().Bits()); err != nil {
		return err
	} else {
		field.SetFloat(floatValue)
		return nil
	}
}

type sliceFieldSetter struct{}

func (s *sliceFieldSetter) Apply(field reflect.Value, value string) error {
	typ := field.Type()
	values := strings.Split(strings.TrimSpace(value), ";")
	count := len(values)
	slice := reflect.MakeSlice(typ, count, count)
	if err := s.fillSlice(slice, values); err != nil {
		return err
	}
	field.Set(slice)
	return nil
}

func (s *sliceFieldSetter) fillSlice(slice reflect.Value, values []string) error {
	if slice.Len() > 0 {
		setter, err := determineFieldSetter(slice.Index(0).Type())
		if err != nil {
			return err
		}
		for i, val := range values {
			elemOf := slice.Index(i)
			setter.Apply(elemOf, val)
		}
	}
	return nil
}

type mapFieldSetter struct{}

func (s *mapFieldSetter) Apply(field reflect.Value, value string) error {
	typ := field.Type()
	mapKeySetter, mapValueSetter, err := s.mapSetters(typ)
	if err != nil {
		return err
	}
	m := reflect.MakeMap(typ)
	keyValues := strings.Split(strings.TrimSpace(value), ";")
	for _, keyValue := range keyValues {
		keyStr, valStr, err := s.splitKeyValuePair(keyValue)
		if err != nil {
			return err
		}
		key := reflect.New(typ.Key()).Elem()
		if err := mapKeySetter.Apply(key, keyStr); err != nil {
			return err
		}
		val := reflect.New(typ.Elem()).Elem()
		if err := mapValueSetter.Apply(val, valStr); err != nil {
			return err
		}
		m.SetMapIndex(key, val)
	}
	field.Set(m)
	return nil
}

func (s *mapFieldSetter) mapSetters(typ reflect.Type) (mapKeySetter FieldSetter, mapValueSetter FieldSetter, err error) {
	mapKeySetter, err = determineFieldSetter(typ.Key())
	if err != nil {
		return
	}
	mapValueSetter, err = determineFieldSetter(typ.Elem())
	if err != nil {
		return
	}
	return
}

func (s *mapFieldSetter) splitKeyValuePair(keyValue string) (key string, value string, err error) {
	splittedPair := strings.Split(keyValue, ":")
	if len(splittedPair) != 2 {
		err = fmt.Errorf("invalid key-pair format is used for map, use key:value format")
		return
	}
	key, value = splittedPair[0], splittedPair[1]
	return
}

type timeDurationFieldSetter struct{}

func (s *timeDurationFieldSetter) Apply(field reflect.Value, value string) error {
	duration, err := time.ParseDuration(value)
	if err != nil {
		return err
	}
	field.SetInt(int64(duration))
	return nil
}
