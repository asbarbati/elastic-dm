package main

import (
	"context"
	"strconv"
	"strings"
	"time"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"github.com/tidwall/gjson"
)

func VerifyConnection(escfg elasticsearch.Config) error {
	logger.Debug("Testing elasticsearch connectivity. Target: " + escfg.Addresses[0])
	var err error

	// Trying to get the ping page
	es, _ := elasticsearch.NewClient(escfg)
	_, err = es.Info()
	if err != nil {
		logger.Error("Error during getting the general info, full error: " + err.Error())
		return err
	}

	return err
}

func VerifyIndex(config elasticsearch.Config, index string) (bool, error) {
	logger.Debug("Verifing the indices '" + index + "' on " + config.Addresses[0])
	// Testing connection
	es, _ := elasticsearch.NewClient(config)
	_, err := es.Info()
	if err != nil {
		logger.Error("Error during connecting.")
		return false, err
	}

	// Trying to getting the document
	res, _ := esapi.IndicesExistsRequest{Index: []string{index}}.Do(context.Background(), es)
	if res.StatusCode == 200 {
		logger.Debug("Indices exists.")
		return true, nil
	}
	logger.Info("Not exists, return: " + strconv.Itoa(res.StatusCode))
	return false, nil
}

func GetDocData(config elasticsearch.Config, index string, docID string) (*strings.Reader, error) {
	logger.Debug("Getting data for" + docID)
	// Testing connection
	es, _ := elasticsearch.NewClient(config)
	_, err := es.Info()
	if err != nil {
		logger.Error("Error during connecting.")
		return nil, err
	}
	doctmp, err := es.Get(index, docID)
	if err != nil {
		return nil, err
	}

	// Clean up data from invalid json
	doc := strings.Split(doctmp.String(), "] ")
	jsoncontent := gjson.Parse(string(doc[1]))
	return strings.NewReader(jsoncontent.Get("_source").String()), nil
}

func GetDocIds(config elasticsearch.Config, mainConfig mainConfigStruct, index string) ([]string, error) {
	logger.Info("Getting Docs IDs from '" + index + "' on " + config.Addresses[0])
	outids := []string{}
	bulkSize := int64(1)

	// Testing connection
	es, _ := elasticsearch.NewClient(config)
	_, err := es.Info()
	if err != nil {
		logger.Error("Error during connecting.")
		return outids, err
	}

	// Get Document number
	stats, err := esapi.IndicesStatsRequest{Index: []string{index}}.Do(context.Background(), es)
	if err != nil {
		return outids, err
	}
	statsjson := strings.Split(stats.String(), "] ")
	jsoncontent := gjson.Parse(statsjson[1])
	totalDocs := jsoncontent.Get("_all.primaries.docs.count").Int()

	// Calculate the bulksize based by the ScrollMultiplier
	if totalDocs >= int64(mainConfig.ScrollMultiplier)*2 {
		bulkSize = totalDocs / int64(mainConfig.ScrollMultiplier)
	}
	logger.Debug("BulkSize: " + strconv.Itoa(int(bulkSize)))

	// Get all the document IDs
	logger.Debug("Getting all the document IDs")
	query := `{ "query": { "match_all": {} }, "_source": false }`
	res, err := es.Search(
		es.Search.WithIndex(index),
		es.Search.WithBody(strings.NewReader(query)),
		es.Search.WithScroll(time.Minute),
		es.Search.WithSize(int(bulkSize)),
	)
	resjson := strings.Split(res.String(), "] ")

	// Parsing the results from the querysearch
	logger.Debug("Parsing the results")
	jsoncontent = gjson.Parse(resjson[1])
	hits := jsoncontent.Get("hits.hits").Array()
	for iter := 0; iter < len(hits); iter++ {
		outids = append(outids, hits[iter].Get("_id").String())
	}
	scrollId := jsoncontent.Get("_scroll_id").String()
	for {
		res, _ := es.Scroll(
			es.Scroll.WithScrollID(scrollId),
			es.Scroll.WithScroll(time.Minute),
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
	logger.Debug("Returning the IDs: " + strings.Join(outids, ","))
	return outids, err
}
