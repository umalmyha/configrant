package configrant

import "github.com/umalmyha/configrant/internal/structs"

// Process apply values to structure fields correspondingly
func Process(from interface{}) error {
	cfg, err := structs.NewParser(from)
	if err != nil {
		return err
	}
	return cfg.MaintainFields()
}
