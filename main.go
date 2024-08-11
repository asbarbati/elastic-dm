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

func main() {
	logger.Info("Process Start")

	configFile := flag.String("config", "./config.json", "Config file to load.")
	debugFlag := flag.Bool("debug", false, "Enable debug mode.")
	modeFlag := flag.String("mode", "sync", "Action to do.")
	flag.Parse()

	if *debugFlag {
		loglevel.Set(slog.LevelDebug)
	}

	// Loading config file
	sharedtimeout, _ := time.ParseDuration("60s")
	var mainConfig mainConfigStruct

	mainConfig, err := loadConfig(*configFile, mainConfig)
	if err != nil {
		return
	}

	esSrcTransportCfg := &http.Transport{
		MaxIdleConnsPerHost:   2,
		ResponseHeaderTimeout: sharedtimeout,
		TLSHandshakeTimeout:   sharedtimeout,
	}
	esDstTransportCfg := &http.Transport{
		MaxIdleConnsPerHost:   2,
		ResponseHeaderTimeout: sharedtimeout,
		TLSHandshakeTimeout:   sharedtimeout,
	}

	if mainConfig.EsSrc.DisableTlsVerify {
		logger.Info("Skipping TLS Verification on the SOURCE Elasticsearch.")
		esSrcTransportCfg.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	}

	if mainConfig.EsDst.DisableTlsVerify {
		logger.Info("Skipping TLS Verification on the DEST Elasticsearch.")
		esDstTransportCfg.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	}

	esSrcCfg := elasticsearch.Config{
		Addresses:                []string{mainConfig.EsSrc.Host},
		Username:                 mainConfig.EsSrc.User,
		Password:                 mainConfig.EsSrc.Passwd,
		CompressRequestBody:      true,
		Transport:                esSrcTransportCfg,
		CompressRequestBodyLevel: 9,
		EnableDebugLogger:        false,
		EnableMetrics:            false,
	}

	logger.Info("Testing Elasticsearch Source target: " + mainConfig.EsSrc.Host)
	err = verifyConnection(esSrcCfg)
	if err != nil {
		return
	}

	esDstCfg := elasticsearch.Config{
		Addresses:           []string{mainConfig.EsDst.Host},
		Username:            mainConfig.EsDst.User,
		Password:            mainConfig.EsDst.Passwd,
		CompressRequestBody: true,
		Transport:           esDstTransportCfg,
	}

	logger.Info("Testing Elasticsearch Dest target: " + mainConfig.EsDst.Host)
	err = verifyConnection(esDstCfg)
	if err != nil {
		return
	}

	logger.Info("Mode: " + *modeFlag)

	if *modeFlag == "sync" {
		err := doSyncMode(esSrcCfg, esDstCfg, mainConfig)
		if err != nil {
			return
		}
	}

}
