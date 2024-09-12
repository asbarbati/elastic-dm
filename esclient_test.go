package main

import (
	"crypto/tls"
	"net/http"
	"testing"

	"github.com/elastic/go-elasticsearch/v8"
)

func TestVerifyConnection(t *testing.T) {
	cfg := elasticsearch.Config{
		Addresses: []string{"https://127.0.0.1:8080"},
	}
	err := VerifyConnection(cfg)
	if err == nil {
		t.Error("Invalid connection.")
	}

	transportcfg := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	cfgvalid := elasticsearch.Config{
		Addresses: []string{"https://127.0.0.1:9200"},
		Username:  "elastic",
		Password:  "elastic",
		Transport: transportcfg,
	}
	err = VerifyConnection(cfgvalid)
	if err != nil {
		t.Error("Connection issue.")
	}
}

func TestVerifyIndex(t *testing.T) {
	transportcfg := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	var status bool

	cfg := elasticsearch.Config{
		Addresses: []string{"https://127.0.0.1:8080"},
	}
	_, err := VerifyIndex(cfg, "elastic-dm-test")
	if err == nil {
		t.Error("Invalid connection not detected.")
	}

	cfg = elasticsearch.Config{
		Addresses: []string{"https://127.0.0.1:9200"},
		Username:  "elastic",
		Password:  "elastic",
		Transport: transportcfg,
	}
	status, err = VerifyIndex(cfg, "elastic-dm-test")
	if err != nil {
		t.Error("Invalid connection.")
	}
	if !status {
		t.Error("Fixture missing.")
	}

	status, _ = VerifyIndex(cfg, "elastic-dm-testnoexists")
	if status {
		t.Error("Index missing exists.")
	}
}

func TestGetDocData(t *testing.T) {
	cfg := elasticsearch.Config{
		Addresses: []string{"https://127.0.0.1:8080"},
	}
	_, err := GetDocData(cfg, "elastic-dm-test", "jTzhrZABVFAGNiUIDqjU")
	if err == nil {
		t.Error("Invalid connection not detected.")
	}

	transportcfg := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	cfg = elasticsearch.Config{
		Addresses: []string{"https://127.0.0.1:9200"},
		Username:  "elastic",
		Password:  "elastic",
		Transport: transportcfg,
	}

	_, err = GetDocData(cfg, "elastic-dm-test", "jTzhrZABVFAGNiUIDqjU")
	if err != nil {
		t.Error("Invalid connection.")
	}
}

func TestGetDocIds(t *testing.T) {
	defatulvalue := SetDefaultValue()
	loadconfig, _ := LoadConfig("tests/assets/config.json", defatulvalue)

	cfg := elasticsearch.Config{
		Addresses: []string{"https://127.0.0.1:8080"},
	}
	_, err := GetDocIds(cfg, loadconfig, "elastic-dm-test")
	if err == nil {
		t.Error("Invalid connection not detected.")
	}

	transportcfg := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	cfg = elasticsearch.Config{
		Addresses: []string{"https://127.0.0.1:9200"},
		Username:  "elastic",
		Password:  "elastic",
		Transport: transportcfg,
	}

	_, err = GetDocIds(cfg, loadconfig, "elastic-dm-test")
	if err != nil {
		t.Error("Invalid connection.")
	}

}
