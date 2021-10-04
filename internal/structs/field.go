package structs

import (
	"os"
	"reflect"
	"strings"
)

type Field struct {
	Elem           reflect.Value
	EnvVarName     string
	DefaultValue   string
	IsConfigurable bool
}

func (f *Field) Set() error {
	value := f.ValueString()
	setter, err := determineFieldSetter(f.Elem.Type())
	if err != nil {
		return err
	}
	return setter.Apply(f.Elem, value)
}

func (f *Field) ValueString() string {
	envVarValue := os.Getenv(f.EnvVarName)
	if envVarValue != "" {
		return envVarValue
	}
	return f.DefaultValue
}

func (f *Field) IsStruct() bool {
	return f.Elem.Kind() == reflect.Struct
}

func NewField(typeOfField reflect.StructField, elemOfField reflect.Value) (field Field) {
	elemOfField = extractFieldElemOf(elemOfField)
	field = Field{
		Elem:           elemOfField,
		IsConfigurable: true,
	}
	tagStr := typeOfField.Tag.Get("cfgrant")
	if !elemOfField.CanSet() || tagStr == "-" {
		field.IsConfigurable = false
		return
	}
	def, env := parseConfigrantTag(tagStr)
	field.DefaultValue = def
	field.EnvVarName = env
	return
}

func extractFieldElemOf(field reflect.Value) (elemOf reflect.Value) {
	elemOf = field
	for elemOf.Kind() == reflect.Ptr {
		if elemOf.IsNil() {
			elemOf.Set(reflect.New(elemOf.Type().Elem()))
		}
		elemOf = elemOf.Elem()
	}
	return
}

func parseConfigrantTag(tagStr string) (def string, env string) {
	if tagStr == "" {
		return
	}
	tagOptions := strings.Split(tagStr, ",")
	for _, tagOption := range tagOptions {
		propValue := strings.SplitN(tagOption, ":", 2)
		prop, value := propValue[0], strings.TrimSpace(propValue[1])
		switch prop {
		case "default":
			def = value
		case "env":
			env = value
		}
	}
	return
}
