package cfgargs

import "strings"

var args = make(map[string]string)

func Lookup(name string) string {
	return args[name]
}

func Parse(arguments []string) {
	for _, arg := range arguments {
		key, value := argKeyValue(arg)
		args[key] = value
	}
}

func argKeyValue(arg string) (key string, val string) {
	argKeyValue := strings.SplitN(arg, "=", 2)
	if len(argKeyValue) == 2 {
		return argKeyValue[0], argKeyValue[1]
	}
	return argKeyValue[0], "true"
}
