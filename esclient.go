package main

import (
	"context"
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"github.com/tidwall/gjson"
)

func verifyConnection(escfg elasticsearch.Config) error {
	logger.Debug("Testing elasticsearch connectivity. Target: " + escfg.Addresses[0])
	var err error
	es, _ := elasticsearch.NewClient(escfg)
	res, err := es.Info()
	if err != nil {
		logger.Error("Error during getting the general info, full error: " + err.Error())
		return err
	}
	if res.StatusCode != 200 {
		logger.Error("Error during getting the general info, full error: " + res.String())
		err = errors.New("returned not 200")
	}
	logger.Debug("Results: " + res.String())

	return err
}

func verifyIndex(config elasticsearch.Config, index string) (bool, error) {
	logger.Debug("Verifing the indices '" + index + "' on " + config.Addresses[0])
	es, err := elasticsearch.NewClient(config)
	if err != nil {
		logger.Error("Error during connecting.")
		return false, err
	}
	res, _ := esapi.IndicesExistsRequest{Index: []string{index}}.Do(context.Background(), es)
	if res.StatusCode == 200 {
		logger.Debug("Indices exists.")
		return true, nil
	}
	logger.Debug("Not exists, return: " + strconv.Itoa(res.StatusCode))
	return false, nil
}

func getDocData(config elasticsearch.Config, index string, docID string) (*strings.Reader, error) {
	logger.Debug("Getting data for" + docID)
	client, err := elasticsearch.NewClient(config)
	if err != nil {
		logger.Error("Error during connecting.")
		return nil, err
	}
	doctmp, err := client.Get(index, docID)
	if err != nil {
		return nil, err
	}

	doc := strings.Split(doctmp.String(), "] ")

	jsondata := cleanupData(doc[1])
	return strings.NewReader(jsondata), nil
}

func cleanupData(jsonin string) string {
	jsoncontent := gjson.Parse(string(jsonin))
	return jsoncontent.Get("_source").String()
}

func getDocIds(config elasticsearch.Config, mainConfig mainConfigStruct, index string) ([]string, error) {
	logger.Info("Getting Docs IDs from '" + index + "' on " + config.Addresses[0])
	outids := []string{}
	bulkSize := int64(1)
	client, err := elasticsearch.NewClient(config)
	if err != nil {
		logger.Error("Error during connecting.")
		return outids, err
	}

	// Get Document number
	stats, err := esapi.IndicesStatsRequest{Index: []string{index}}.Do(context.Background(), client)
	if err != nil {
		return outids, err
	}
	statsjson := strings.Split(stats.String(), "] ")
	jsoncontent := gjson.Parse(statsjson[1])
	totalDocs := jsoncontent.Get("_all.primaries.docs.count").Int()

	if totalDocs >= int64(mainConfig.ScrollMultiplier)*2 {
		bulkSize = totalDocs / int64(mainConfig.ScrollMultiplier)
	}

	query := `{ "query": { "match_all": {} }, "_source": false }`
	res, err := client.Search(
		client.Search.WithIndex(index),
		client.Search.WithBody(strings.NewReader(query)),
		client.Search.WithScroll(time.Minute),
		client.Search.WithSize(int(bulkSize)),
	)

	resjson := strings.Split(res.String(), "] ")

	jsoncontent = gjson.Parse(resjson[1])
	hits := jsoncontent.Get("hits.hits").Array()
	for iter := 0; iter < len(hits); iter++ {
		outids = append(outids, hits[iter].Get("_id").String())
	}
	scrollId := jsoncontent.Get("_scroll_id").String()
	for {
		res, _ := client.Scroll(
			client.Scroll.WithScrollID(scrollId),
			client.Scroll.WithScroll(time.Minute),
		)
		resjson := strings.Split(res.String(), "] ")
		jsoncontent = gjson.Parse(resjson[1])
		hits := jsoncontent.Get("hits.hits").Array()

		if len(hits) == 0 {
			break
		}

		for iter := 0; iter < len(hits); iter++ {
			outids = append(outids, hits[iter].Get("_id").String())
		}
	}
	return outids, err
}
