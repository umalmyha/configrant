package configrant

import (
	"os"

	"github.com/umalmyha/configrant/internal/cfgargs"
	"github.com/umalmyha/configrant/internal/structs"
)

// Process apply values to structure fields correspondingly
func Process(from interface{}) error {
	cfgargs.Parse(os.Args)
	cfg, err := structs.NewParser(from)
	if err != nil {
		return err
	}
	return cfg.MaintainFields()
}
