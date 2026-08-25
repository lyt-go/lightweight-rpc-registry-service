package service

import (
	"fmt"

	"lightweightrpc/internal/config"
)

func PlanCodec(defaults config.CodecDefaults, name string) error {
	if defaults.Validator != nil && !defaults.Validator.Allow(name) {
		return fmt.Errorf("codec name rejected: %s", name)
	}
	defaults.Aliases[name] = "ready"
	return nil
}
