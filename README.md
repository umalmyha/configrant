# configrant

Package configrant is used to maintain configuration structures with command line arguments, environment variables and default values by tagging structure fields.

## Installation

```
go get github.com/umalmyha/configrant
```

## Usage

Before usage, please, see detailed [documentation](https://pkg.go.dev/github.com/umalmyha/configrant). There are some edge cases you should be aware of.

See below example of processing tagged configuration strucutre:

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
	Timeout   time.Duration  `cfgrant:"default:5s,arg:--timeout"`
	Substruct ConfigSubstruct
}

cfg := &Config{}
err := configrant.Process(cfg)
if err != nil {
	fmt.Println(err.Error())
}
```

## Contribution and bugs

Any ideas for package improvement and contribution are appreciated.
