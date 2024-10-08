package main

import (
	"errors"
	"os"
	"strings"

	"github.com/tidwall/gjson"
)

func ReadConfig(jsondata gjson.Result, mainConfig mainConfigStruct, defaultvalue defaultValueStruct) (mainConfigStruct, error) {
	var err error
	mandatoryFields := []string{
		"es_src", "es_src.host", "es_src.passwd", "es_src.user", "es_src.indices",
		"es_dst", "es_dst.host", "es_dst.passwd", "es_dst.user", "es_dst.indices",
		"mode",
	}

	// Check mandatory elements
	logger.Info("Verify the mandatory fields...")
	for iter := 0; iter < len(mandatoryFields); iter++ {
		logger.Debug("Verifing field '" + mandatoryFields[iter] + "' exists...")
		if !jsondata.Get(mandatoryFields[iter]).Exists() {
			logger.Error("The field '" + mandatoryFields[iter] + "' its mandatory. Exiting")
			err = errors.New("field_missing")
			break
		}
	}

	// Catch if some field missing
	if err != nil {
		return mainConfig, err
	}

	logger.Info("Loading Config and defaults...")

	// ScrollMultiplier
	if jsondata.Get("scrollmultiplier").Exists() {
		mainConfig.ScrollMultiplier = int(jsondata.Get("scrollmultiplier").Int())
	} else {
		mainConfig.ScrollMultiplier = defaultvalue.ScrollMultiplier
	}

	// Worker details
	if jsondata.Get("workers").Exists() {
		mainConfig.Workers = int(jsondata.Get("workers").Int())
	} else {
		mainConfig.Workers = defaultvalue.Workers
	}

	// Daemon mode
	if jsondata.Get("mode").Exists() {
		mainConfig.Mode = jsondata.Get("mode").String()
	} else {
		mainConfig.Mode = defaultvalue.Mode
	}

	// Source Elasticsearch Config
	mainConfig.EsSrc.Host = jsondata.Get("es_src.host").String()
	mainConfig.EsSrc.User = jsondata.Get("es_src.user").String()
	mainConfig.EsSrc.Passwd = jsondata.Get("es_src.passwd").String()
	mainConfig.EsSrc.Indices = strings.Split(jsondata.Get("es_src.indices").String(), ",")
	if jsondata.Get("es_src.disabletlsverify").Exists() {
		mainConfig.EsSrc.DisableTlsVerify = jsondata.Get("es_src.disabletlsverify").Bool()
	} else {
		mainConfig.EsSrc.DisableTlsVerify = false
	}

	// Destination Elasticsearch Config
	mainConfig.EsDst.Host = jsondata.Get("es_dst.host").String()
	mainConfig.EsDst.User = jsondata.Get("es_dst.user").String()
	mainConfig.EsDst.Passwd = jsondata.Get("es_dst.passwd").String()
	mainConfig.EsDst.Indices = strings.Split(jsondata.Get("es_dst.indices").String(), ",")
	if jsondata.Get("es_dst.disabletlsverify").Exists() {
		mainConfig.EsDst.DisableTlsVerify = jsondata.Get("es_dst.disabletlsverify").Bool()
	} else {
		mainConfig.EsDst.DisableTlsVerify = false
	}

	return mainConfig, err
}

func LoadConfig(fpath string, defaultvalues defaultValueStruct) (mainConfigStruct, error) {
	var mainConfig mainConfigStruct
	// Check if file exists
	logger.Debug("Checking config file...")
	_, err := os.Stat(fpath)
	if errors.Is(err, os.ErrNotExist) {
		logger.Error("Config file '" + fpath + "' not found. Please checkout the README.")
		return mainConfig, err
	}

	// Try to open it
	logger.Debug("Opening config file...")
	configFile, err := os.ReadFile(fpath)
	if err != nil {
		logger.Error("Error opening the config file. Exiting.")
		return mainConfig, err
	}

	// Load the content
	configcontent := string(configFile)
	logger.Debug("Decoding the config file as json...")
	if !gjson.Valid(configcontent) {
		logger.Error("Invalid json, exiting.")
		return mainConfig, errors.New("")
	}

	jsoncontent := gjson.Parse(configcontent)
	mainConfig, err = ReadConfig(jsoncontent, mainConfig, defaultvalues)
	return mainConfig, err
}
