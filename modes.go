package main

import (
	"context"
	"errors"
	"strconv"
	"strings"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esutil"
)

func doSyncMode(srccfg elasticsearch.Config, dstcfg elasticsearch.Config, mainconfig mainConfigStruct) error {
	for iter := 0; iter < len(mainconfig.EsSrc.Indices); iter++ {
		var dstIndex string
		var dstIndexExist bool
		var srcDocIDs []string
		var dstDocIDs []string
		var diffDocIDs []string

		// Verify if the number of the index its the same
		if len(mainconfig.EsSrc.Indices) != len(mainconfig.EsDst.Indices) {
			logger.Error("SrcIndices and DestIndices don't have the same number of indexes. Use #COPY# as placeholder for copy the name as source.")
			return errors.New("slice not matched")
		}

		// Verify if the SrcIndex Exists. Else skip.
		srcExist, err := verifyIndex(srccfg, mainconfig.EsSrc.Indices[iter])
		if err != nil {
			return err
		}
		if !srcExist {
			logger.Info("The indices '" + mainconfig.EsSrc.Indices[iter] + "' not exists. Skip")
			continue
		}

		// Check the dest index
		if mainconfig.EsDst.Indices[iter] == "#COPY#" {
			dstIndex = mainconfig.EsSrc.Indices[iter]
		} else {
			dstIndex = mainconfig.EsDst.Indices[iter]
		}

		// Verify if dest index exist
		dstIndexExist, err = verifyIndex(dstcfg, dstIndex)
		if err != nil {
			return err
		}

		logger.Debug("The dest indices '" + dstIndex + "' returns: " + strconv.FormatBool(dstIndexExist))

		logger.Info("Getting IDs from the source index '" + mainconfig.EsSrc.Indices[iter])
		srcDocIDs, _ = getDocIds(srccfg, mainconfig, mainconfig.EsSrc.Indices[iter])

		if dstIndexExist {
			dstDocIDs, _ = getDocIds(dstcfg, mainconfig, dstIndex)
			if len(dstDocIDs) != 0 {
				for srcdociter := 0; srcdociter < len(srcDocIDs); srcdociter++ {
					docExist := false
					for destdociter := 0; destdociter < len(dstDocIDs); destdociter++ {
						if srcDocIDs[srcdociter] == dstDocIDs[destdociter] {
							docExist = true
							break
						}
					}
					if !docExist {
						diffDocIDs = append(diffDocIDs, srcDocIDs[srcdociter])
					}
				}
			}
		} else {
			diffDocIDs = append(diffDocIDs, srcDocIDs...)
		}

		logger.Debug(strings.Join(diffDocIDs, ","))
		logger.Info("Prepare the ES Bulk indexer...")
		es, _ := elasticsearch.NewClient(srccfg)
		esIndexer, err := esutil.NewBulkIndexer(esutil.BulkIndexerConfig{
			Client:     es,
			Index:      dstIndex,
			NumWorkers: 40,
			FlushBytes: 1e+6,
		})

		if err != nil {
			logger.Error(err.Error())
			return err
		}

		for diffDocIDIter := 0; diffDocIDIter < len(diffDocIDs); diffDocIDIter++ {
			docdata, err := getDocData(srccfg, mainconfig.EsSrc.Indices[iter], diffDocIDs[diffDocIDIter])
			if err != nil {
				return err
			}
			err = esIndexer.Add(
				context.Background(),
				esutil.BulkIndexerItem{
					Action:     "index",
					DocumentID: diffDocIDs[diffDocIDIter],
					Body:       docdata,
				},
			)
			if err != nil {
				logger.Error(err.Error())
				return err
			}
		}

		if err := esIndexer.Close(context.Background()); err != nil {
			logger.Error("Error")
		}
		logger.Info("Done.")

	}
	return nil
}
