package main

import (
	"crypto/tls"
	"flag"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/elastic/go-elasticsearch/v8"
)

var logfile, _ = os.OpenFile("./elastic-dm.json.log", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
var loglevel = new(slog.LevelVar)
var logger = slog.New(slog.NewJSONHandler(logfile, &slog.HandlerOptions{Level: loglevel}))

func SetDefaultValue() defaultValueStruct {
	var defaultsetting defaultValueStruct

	// Shared TTL for elasticsearch for transport
	defaultsetting.ESSharedTimeout, _ = time.ParseDuration("60s")

	// Goroutine workers to use
	defaultsetting.Workers = 10

	// Multiplier to use during the bulk requests on ES
	defaultsetting.ScrollMultiplier = 10

	// Default mode if not specified
	defaultsetting.Mode = "None"

	return defaultsetting
}

func StartProcess(mainConfig mainConfigStruct, defaultsettings defaultValueStruct) error {
	// Set up the shared timeouts
	esTransportCfg := &http.Transport{
		ResponseHeaderTimeout: defaultsettings.ESSharedTimeout,
		TLSHandshakeTimeout:   defaultsettings.ESSharedTimeout,
	}

	// Disable the TLSVerification if requested
	if mainConfig.EsSrc.DisableTlsVerify {
		logger.Info("Skipping TLS Verification on the SOURCE Elasticsearch.")
		esTransportCfg.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	}
	if mainConfig.EsDst.DisableTlsVerify {
		logger.Info("Skipping TLS Verification on the DEST Elasticsearch.")
		esTransportCfg.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	}

	// Loading the configuration for Source Elastic
	esSrcCfg := elasticsearch.Config{
		Addresses:                []string{mainConfig.EsSrc.Host},
		Username:                 mainConfig.EsSrc.User,
		Password:                 mainConfig.EsSrc.Passwd,
		CompressRequestBody:      true,
		Transport:                esTransportCfg,
		CompressRequestBodyLevel: 9,
		EnableDebugLogger:        false,
		EnableMetrics:            false,
	}

	// Verify if Elastic Source connection is valid
	logger.Info("Testing Elasticsearch Source target: " + mainConfig.EsSrc.Host)
	err := VerifyConnection(esSrcCfg)
	if err != nil {
		return err
	}

	// Loading the configuration for Dest Elastic
	esDstCfg := elasticsearch.Config{
		Addresses:           []string{mainConfig.EsDst.Host},
		Username:            mainConfig.EsDst.User,
		Password:            mainConfig.EsDst.Passwd,
		CompressRequestBody: true,
		Transport:           esTransportCfg,
	}

	// Verify if Elastic Dest connection is valid
	logger.Info("Testing Elasticsearch Dest target: " + mainConfig.EsDst.Host)
	err = VerifyConnection(esDstCfg)
	if err != nil {
		return err
	}

	// Doing stuff based by the Mode
	// Sync = Syncronize the source to the dest index
	logger.Info("Mode: " + mainConfig.Mode)
	if mainConfig.Mode == "sync" {
		err := DoSyncMode(esSrcCfg, esDstCfg, mainConfig)
		if err != nil {
			return err
		}
	}
	return nil
}

func main() {
	logger.Info("Process Start")

	// Argument parse
	configFile := flag.String("config", "./config.json", "Config file to load.")
	debugFlag := flag.Bool("debug", false, "Enable debug mode.")
	flag.Parse()

	// Set debug log if specified
	if *debugFlag {
		loglevel.Set(slog.LevelDebug)
	}

	// Loading the configuration
	defaultsettings := SetDefaultValue()
	mainConfig, err := LoadConfig(*configFile, defaultsettings)

	if err != nil {
		logger.Error("Error during loading the config file, please check.")
		return
	}

	StartProcess(mainConfig, defaultsettings)
	logger.Info("Process End.")
}
