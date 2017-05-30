package config

import (
	"io/ioutil"

	"github.com/alvinfeng/mosaic/storage/driver/filesystem"
	yaml "gopkg.in/yaml.v2"
)

type Config struct {
	Addr        string            `yaml:"addr"`
	Insecure    bool              `yaml:"insecure"`
	CertPath    string            `yaml:"certpath"`
	KeyPath     string            `yaml:"keypath"`
	StorageType string            `yaml:"storagetype"`
	Fs          filesystem.Config `yaml:"filesystem"`
	CacheType   string            `yaml:"cachetype"`
}

// LoadConfig loads a Config from a yaml file
func LoadConfig(filename string) (*Config, error) {
	contents, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	c := &Config{}
	if err = yaml.Unmarshal(contents, c); err != nil {
		return nil, err
	}

	return c, nil
}
