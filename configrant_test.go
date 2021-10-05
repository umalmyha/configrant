package configrant

import (
	"reflect"
	"testing"

	"github.com/umalmyha/configrant/internal/structs"
)

type ConfigSubstruct struct {
	Subname string  `cfgrant:"env:SUBNAME_ENV,default:SubConfig"`
	Percent float32 `cfgrant:"default:3.32"`
}

type Config struct {
	//lint:ignore U1000 we must test that unexportable field is ignored even if tagged
	private   string         `cfgrant:"default:private"`
	Name      string         `cfgrant:"env:NAME_ENV"`
	Url       string         `cfgrant:"default:http://localhost:3000"`
	Retries   int            `cfgrant:"env:RETRIES_ENV,default:3"`
	OwnerPtr  *string        `cfgrant:"env:OWNER,default:James"`
	Pass_hash string         `cfgrant:"-"`
	Bytes     []byte         `cfgrant:"default:1;2;3;4;5"`
	Sequence  map[string]int `cfgrant:"default:second:2;third:3;first:1"`
	IsAsync   bool           `cfgrant:"default:true"`
	// Timeout   time.Duration  `cfgrant:"default:5s"` // TODO: add support as it used pretty frequently in config
	Password  string
	Substruct ConfigSubstruct
}

func TestProcessSuccess(t *testing.T) {
	cfg := &Config{
		Password:  "password",
		Pass_hash: "d1e8a70b5ccab1dc2f56bbf7e99f064a660c08e361a35751b9c483c88943d082",
	}

	// set environment variables
	t.Setenv("NAME_ENV", "configuration")
	t.Setenv("OWNER", "Ronald")

	t.Log("Expect success parse of complex config struct")

	// error shouldn't occur
	if err := Process(cfg); err != nil {
		t.Fatalf("Error occured during parsing %s", err.Error())
	}
	// unexported field is not settable, so should stay zero-valued
	if cfg.private != "" {
		t.Errorf(`Expect field 'private' to be equal "", got %s`, cfg.private)
	}

	// NAME_ENV is defined, so value must be taken from env variable
	if cfg.Name != "configuration" {
		t.Errorf(`Expect field 'Name' to be equal "configuration", got %s`, cfg.Name)
	}

	// no env property specified, so expect defaul value
	if cfg.Url != "http://localhost:3000" {
		t.Errorf(`Expect field 'Url' to be equal "http://localhost:3000", got %s`, cfg.Url)
	}

	// env is specified, but it is not set, so default value is expected
	if cfg.Retries != 3 {
		t.Errorf(`Expect field 'Retries' to be equal 3, got %d`, cfg.Retries)
	}

	// pointer shouldn't cause any troubles and field must be initialized properly
	// env is specified as well as default, but env has priority
	if cfg.OwnerPtr == nil {
		t.Error("Expect field 'OwnerPtr' to be not nil")
	} else if *cfg.OwnerPtr != "Ronald" {
		t.Errorf(`Expect field 'OwnerPtr' to be equal "Ronald", got %s`, *cfg.OwnerPtr)
	}

	// field is ignored with "-", so value should stay unchanged
	if cfg.Pass_hash != "d1e8a70b5ccab1dc2f56bbf7e99f064a660c08e361a35751b9c483c88943d082" {
		t.Errorf(`Expect field 'Pass_hash' to be equal "d1e8a70b5ccab1dc2f56bbf7e99f064a660c08e361a35751b9c483c88943d082", got %s`, cfg.Pass_hash)
	}

	// slice must be initialized
	// default value must be taken
	if cfg.Bytes == nil {
		t.Error("Expect field 'Bytes' to be not nil")
	} else if !reflect.DeepEqual(cfg.Bytes, []byte{1, 2, 3, 4, 5}) {
		t.Errorf("Expect field 'Bytes' to be equal [1 2 3 4 5], got %v", cfg.Bytes)
	}

	// map must be initialized
	// default value must be taken
	if cfg.Sequence == nil {
		t.Error("Expect field 'Sequence' to be not nil")
	} else if !reflect.DeepEqual(cfg.Sequence, map[string]int{"second": 2, "third": 3, "first": 1}) {
		t.Errorf("Expect field 'Sequence' to be equal map[first:1 second:2 third:3], got %v", cfg.Sequence)
	}

	// default true value is expected
	if cfg.IsAsync != true {
		t.Errorf("Expect field 'isAsync' to be equal true, got %t", cfg.IsAsync)
	}

	// Field is not tagged, but not ignored and set empty value instead of defined on struct initialization
	if cfg.Password != "" {
		t.Errorf(`Expect field 'Password' to be equal "", got %s`, cfg.Password)
	}

	// Expect to be default, because env is specified, but not set
	if cfg.Substruct.Subname != "SubConfig" {
		t.Errorf(`Expect inner struct field 'Subname' to be equal "SubConfig", got %s`, cfg.Substruct.Subname)
	}

	if cfg.Substruct.Percent != 3.32 {
		t.Errorf("Expect inner struct field 'Percent' to be equal 3.32, got %.2f", cfg.Substruct.Percent)
	}
}

func TestProcessNotStructPassed(t *testing.T) {
	t.Log("Try passing non-struct values")
	if err := Process(5); err == nil {
		t.Errorf("Non-struct value 5 is passed and got no error")
	}
	if err := Process("Hello World"); err == nil {
		t.Errorf(`Non-struct value "Hello World" is passed and got no error`)
	}
	if err := Process(3.32); err == nil {
		t.Errorf("Non-struct value 3.32 is passed and got no error")
	}
	if err := Process(true); err == nil {
		t.Errorf("Non-struct value true is passed and got no error")
	}
}

func TestProcessNotPtrStructPassed(t *testing.T) {
	t.Log("Try passing not struct pointer")
	cfg := Config{}
	if err := Process(cfg); err == nil || err != structs.ErrNotPtrStruct {
		t.Fatal("Not pointer to struct has been passed and was no error")
	}
}
