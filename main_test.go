package main

import "testing"

func TestSetDefaultValue(t *testing.T) {
	defaultsettings := SetDefaultValue()
	if defaultsettings.Workers != 10 {
		t.Error("Worker not default.")
	}
}

func TestStartProcess(t *testing.T) {
	defaultsettings := SetDefaultValue()
	mainConfig, err := LoadConfig("tests/assets/config.json", defaultsettings)
	if err != nil {
		t.Error("Error loading config")
	}

	err = StartProcess(mainConfig, defaultsettings)
	if err != nil {
		t.Error("Error during starting process")
	}
}
