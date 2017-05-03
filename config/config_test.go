package config

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type ConfigTestSuite struct {
	suite.Suite
	config Config
}

func (suite *ConfigTestSuite) SetupTest() {
	suite.config = Config{
		StorageType: "filesystem",
	}
}

func (suite *ConfigTestSuite) TestLoadConfig() {
	c, err := LoadConfig("./test.yaml")
	suite.Nil(err)
	suite.Equal(*c, suite.config)
}

func TestConfigTestSuite(t *testing.T) {
	suite.Run(t, new(ConfigTestSuite))
}
