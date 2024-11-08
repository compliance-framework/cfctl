package create

import (
	"cfctl/common"
	"testing"
)

func setupExpectedConfig() common.CLIConfig {
	expectedContext := make(map[string]common.Context)
	expectedContext["test"] = common.Context{URL: "test"}
	expectedConfig := common.CLIConfig{Contexts: expectedContext}
	return expectedConfig
}

func TestCreateContextFromEmpty(t *testing.T) {
	testConfig := common.CLIConfig{}

	expectedConfig := setupExpectedConfig()

	configWithContext := CreateContext(testConfig, "test", "test")

	if configWithContext.Contexts["test"] != expectedConfig.Contexts["test"] {
		t.Errorf("unexpected result: got %v but expected %v", configWithContext.Contexts["test"], expectedConfig.Contexts["test"])
	}
}

func TestCreateContextFromExisting(t *testing.T) {
	testContextMap := make(map[string]common.Context)
	testContextMap["test"] = common.Context{URL: "overwrite"}
	testConfig := common.CLIConfig{Default: "test", Contexts: testContextMap}

	expectedConfig := setupExpectedConfig()

	configWithContext := CreateContext(testConfig, "test", "test")

	if configWithContext.Contexts["test"] != expectedConfig.Contexts["test"] {
		t.Errorf("unexpected result: got %v but expected %v", configWithContext.Contexts["test"], expectedConfig.Contexts["test"])
	}
}
