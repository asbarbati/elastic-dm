package main

import (
	"crypto/tls"
	"net/http"
	"testing"

	"github.com/elastic/go-elasticsearch/v8"
)

func TestDoSyncMode(t *testing.T) {
	defatulvalue := SetDefaultValue()
	loadconfig, _ := LoadConfig("tests/assets/config.json", defatulvalue)

	transportcfg := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	cfgsrc := elasticsearch.Config{
		Addresses: []string{"https://127.0.0.1:9200"},
		Username:  "elastic",
		Password:  "elastic",
		Transport: transportcfg,
	}
	cfgdst := elasticsearch.Config{
		Addresses: []string{"https://127.0.0.1:9200"},
		Username:  "elastic",
		Password:  "elastic",
		Transport: transportcfg,
	}

	err := DoSyncMode(cfgsrc, cfgdst, loadconfig)
	if err != nil {
		t.Error("Error first data sync")
	}
	// Double run for matching the resume tests
	err = DoSyncMode(cfgsrc, cfgdst, loadconfig)
	if err != nil {
		t.Error("Error first data sync")
	}
}
