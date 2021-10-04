package structs

import (
	"errors"
	"reflect"
)

var ErrNotPtrStruct = errors.New("configuration must be a pointer to a struct")

type Parser struct {
	ValueOf reflect.Value
	ElemOf  reflect.Value
	TypeOf  reflect.Type
}

func NewParser(from interface{}) (Parser, error) {
	var cfg Parser
	valueOf := reflect.ValueOf(from)
	if valueOf.Kind() != reflect.Ptr {
		return cfg, ErrNotPtrStruct
	}
	elemOf := valueOf.Elem()
	if elemOf.Kind() != reflect.Struct {
		return cfg, ErrNotPtrStruct
	}
	typeOf := elemOf.Type()
	cfg = Parser{
		valueOf,
		elemOf,
		typeOf,
	}
	return cfg, nil
}

func (cfg Parser) MaintainFields() error {
	fields, err := cfg.collectConfigFields()
	if err != nil {
		return err
	}
	for _, field := range fields {
		field.Set()
	}
	return nil
}

func (cfg Parser) collectConfigFields() ([]Field, error) {
	fields := make([]Field, 0)
	for i := 0; i < cfg.ElemOf.NumField(); i++ {
		field := NewField(cfg.TypeOf.Field(i), cfg.ElemOf.Field(i))
		switch {
		case !field.IsConfigurable:
			continue
		case field.IsStruct():
			subcfg, err := NewParser(field.Elem.Addr().Interface())
			if err != nil {
				return nil, err
			}
			substructureFields, err := subcfg.collectConfigFields()
			if err != nil {
				return nil, err
			}
			fields = append(fields, substructureFields...)
		default:
			fields = append(fields, field)
		}
	}
	return fields, nil
}
