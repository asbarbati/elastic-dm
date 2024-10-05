#!/bin/bash

READY=0
until [[ ${READY} == "1" ]]; do
    curl -X GET https://127.0.0.1:9200 -k -u elastic:elastic > /dev/null 2>&1 
    if [[ $? == "0" ]]; then
        READY="1"
    fi
    sleep 1
done

echo "Elasticsearch Ready."

curl -u elastic:elastic -k -X POST "https://localhost:9200/_bulk?pretty" -H 'Content-Type: application/json' -d'
{"index" : {"_index" : "elastic-dm-test", "_id" : "jTzhrZABVFAGNiUIDqjU" } }
{"name":"Foobar","last":"Bar","address":"","ip":"1.1.1.1","email":"info@example.com","epoch":1582889588.3967564,"words":["property","truth","across"]}
{"index" : {"_index" : "elastic-dm-test", "_id" : "abcde" } }
{"name":"Foo","last":"Bar","address":"","ip":"1.1.1.1","email":"info@example.com","epoch":1582889588.3967564,"words":["property","truth","across"]}
'
