package config

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

var ErrUnsupportedFormat = errors.New("formato de arquivo não suportado")

type Config struct {
	Keywords []string `yaml:"keywords" json:"keywords"`
	Domains  []string `yaml:"domains"  json:"domains"`
}

func Load(path string) (*Config, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	cfg := &Config{}
	ext := strings.ToLower(filepath.Ext(path))
	switch ext {
	case ".yaml", ".yml":
		if err := yaml.Unmarshal(data, cfg); err != nil {
			return nil, err
		}
	case ".json":
		if err := json.Unmarshal(data, cfg); err != nil {
			return nil, err
		}
	default:
		return nil, ErrUnsupportedFormat
	}
	return cfg, nil
}
