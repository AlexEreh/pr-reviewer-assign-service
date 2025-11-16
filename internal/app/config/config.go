package config

import (
	"errors"
	"fmt"

	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
)

func FromFile(path string) (*koanf.Koanf, error) {
	k := koanf.New(".")

	err := k.Load(file.Provider(path), yaml.Parser())
	if err != nil {
		return nil, errors.Join(fmt.Errorf("unable to parse config: %w", err))
	}

	return k, nil
}
