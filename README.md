# configrant

Package configrant is used to maintain configuration structures with environment variables and default values by simply tagging structure fields.

## Installation

```
go get github.com/umalmyha/configrant
```

## Configrant tag

Configrant uses special `cfgrant` tag to provide details on field maintenance.

```go
type Config struct {
	Url 	string `cfgrant:"env:ENV_URL,default:http://localhost:8080"`
	IsAsync	bool   `cfgrant:"default:true"`
}
```

Two options are supported:

```
env     - environment variable name
default - default value
```

For struct example mentioned above, we tell configrant:

1. Take value for Url field from environment variable ENV_URL and, if this variable is not defined (or initial), then apply default value.

2. Take default value for field isAsync.

So, environment variable always has priority over default value.

## Simple example

Please, see below simple usage of configrant module:

```go
type Config struct {
	Url 	string `cfgrant:"env:ENV_URL,default:http://localhost:8080"`
	IsAsync	bool   `cfgrant:"default:true"`
}

cfg := &Config{}
err := configrant.Process(cfg)
if err != nil {
	fmt.Println(err.Error())
}
```

That's it. After execution of this code your configuration struct will be maintained with values accordingly.

## How to tag

There are some rules which you should follow, so your struct will be maintained correctly.
Under the hood configrant uses standard reflect package, so some rules coming from here, more precisely from [The Laws of Reflection](https://go.dev/blog/laws-of-reflection)

You must pass pointer to a configuration struct, so it will be settable:

```go
// Correct usage -> pointer to a struct is passed
cfg := &Config{}
err := configrant.Process(cfg)

// Incorrect usage -> value semantic is used
cfg := Config{}
err := configrant.Process(cfg)
```

You must define exported fields (capitalized) in your configuration structures, so it will be settable:

```go
type Config struct {
	url 	string `cfgrant:"env:ENV_URL,default:http://localhost:8080"` // field is unexportable, so not settable -> will be ignored even if tagged
	IsAsync	bool   `cfgrant:"default:true"` // field is exportable, so settable
}
```

For tag options, please, use comma as a separator. Using different separator might cause unexpected behaviour. Adding non-existing option will be ignored:

```go
type Config struct {
	Url 	string `cfgrant:"env:ENV_URL,default:http://localhost:8080"` // correct tag
	IsAsync	bool   `cfgrant:"env:ENV_ASYNC;default:true"` // incorrect tag -> ';' delimiter is used instead of ','
	Retries int    `cfgrant:"env:ENV_RETRIES,option:value"` // correct tag, but property 'option' is ignored
}
```

All basic Go types are supported. There is also support for time.Duration as this type is used pretty frequently for configuration (timeout, etc.). Please, see the whole list:

- string
- bool
- int
- float
- slice
- map
- time.Duration

Slice elements must be separated by semicolon:

```go
type Config struct {
	Sl1 []string `cfgrant:"default:a;b;c;d;e"` // slice elements are defined correctly
	Sl2 []string `cfgrant:"default:a b c d e"` // slice elements are defined incorrectly -> space separator is used
}
```

Map elements must be separated by semicolon and each key-value pair must be defained in format `key:value`:

```go
type Config struct {
	M1 map[string]int `cfgrant:"default:a:1;b:2;c:3"` // map elements are defined correctly
	M2 map[string]int `cfgrant:"default:a-1;b-2;c-3"` // map elements are defined incorrectly -> key-value format is incorrect
}
```

Embedded strucutres are supported as well. Tag is not required for them:

```go
type SubConfig struct {
	Name string `cfgrant:"env:ENV_NAME"`
}
type Config struct {
	IsAsync	bool `cfgrant:"env:ENV_ASYNC,default:true"`
	Sc SubConfig
}
```

You must set fields for ignoring explicitly by adding `cfgrant:"-"` or values you specified on structure initialization will be overwritten:

```go
type Config struct {
	Count	int `cfgrant:"-"` // ingored -> if on initialization we set Count equal to 5 it will stay unchanged
	Retries int	// not tagged and not ignored -> if on initialization we set Retries equal to 3 it will be set to 0 (zero value)
}
```

You can use pointers as well:

```go
type Config struct {
	Name *string `cfgrant:"default:James"`
}
```

## Compex example

Please, see below some complex example with different field types and embedded structure:

```go
type ConfigSubstruct struct {
	Subname string  `cfgrant:"env:SUBNAME_ENV,default:SubConfig"`
	Percent float32 `cfgrant:"default:3.32"`
}

type Config struct {
	Url       string         `cfgrant:"default:http://localhost:3000"`
	Retries   int            `cfgrant:"env:RETRIES_ENV,default:3"`
	PassHash  string         `cfgrant:"-"`
	Bytes     []byte         `cfgrant:"default:1;2;3;4;5"`
	Sequence  map[string]int `cfgrant:"default:second:2;third:3;first:1"`
	IsAsync   bool           `cfgrant:"default:true"`
	Timeout   time.Duration  `cfgrant:"default:5s"`
	Substruct ConfigSubstruct
}

cfg := &Config{}
err := configrant.Process(cfg)
if err != nil {
	fmt.Println(err.Error())
}
```

## Contribution and bugs

Any ideas for package improvement and contribution is appreciated.
