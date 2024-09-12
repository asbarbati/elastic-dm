package main

import (
	"os"
	"testing"

	"github.com/tidwall/gjson"
)

func TestReadConfig(t *testing.T) {
	var mainConfig mainConfigStruct
	defaultvalues := SetDefaultValue()
	configFile, _ := os.ReadFile("tests/assets/config.json")
	jsoncontent := gjson.Parse(string(configFile))
	mainConfig, err := ReadConfig(jsoncontent, mainConfig, defaultvalues)
	if err != nil {
		t.Error(err)
	}
	if mainConfig.Mode != "sync" {
		t.Error("No action detected.")
	}

	configFile, _ = os.ReadFile("tests/assets/config_nofield.json")
	jsoncontent = gjson.Parse(string(configFile))
	mainConfig, err = ReadConfig(jsoncontent, mainConfig, defaultvalues)
	if err == nil {
		t.Error("No field mandatory detected.")
	}

	configFile, _ = os.ReadFile("tests/assets/config_nodefault.json")
	jsoncontent = gjson.Parse(string(configFile))
	mainConfig, _ = ReadConfig(jsoncontent, mainConfig, defaultvalues)
	if mainConfig.Workers != 10 {
		t.Error("No default catched")
	}

}

func TestLoadConfig(t *testing.T) {
	defatulvalue := SetDefaultValue()
	_, err := LoadConfig("tests/assets/config_invalid.json", defatulvalue)
	if err == nil {
		t.Error("No json invalid detected.")
	}

	_, err = LoadConfig("tests/assets/config_noperm.json", defatulvalue)
	if err == nil {
		t.Error("No permission access detected.")
	}

	loadconfig, err := LoadConfig("tests/assets/config.json", defatulvalue)
	if err != nil {
		t.Error(err)
	}
	if loadconfig.Workers != 20 {
		t.Error("No Worker detected")
	}
}
