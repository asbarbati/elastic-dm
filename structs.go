package main

import "time"

type defaultValueStruct struct {
	ESSharedTimeout  time.Duration
	Workers          int
	ScrollMultiplier int
	Mode             string
}

type esConfigStruct struct {
	Host             string   `json:"host"`
	Tls              bool     `json:"tls"`
	DisableTlsVerify bool     `json:"disabletlsverify"`
	Port             int      `json:"port"`
	User             string   `json:"user"`
	Passwd           string   `json:"passwd"`
	Indices          []string `json:"indices"`
}

type mainConfigStruct struct {
	EsSrc            esConfigStruct `json:"es_src"`
	EsDst            esConfigStruct `json:"es_dst"`
	Mode             string         `json:"mode"`
	ScrollMultiplier int            `json:"scrollmultiplier"`
	Workers          int            `json:"workers"`
}
