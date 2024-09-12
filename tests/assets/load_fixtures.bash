#!/bin/bash

curl -u elastic:elastic -k -X POST "https://localhost:9200/_bulk?pretty" -H 'Content-Type: application/json' -d'
{"index" : {"_index" : "elastic-dm-test", "_id" : "jTzhrZABVFAGNiUIDqjU" } }
{"name":"Foobar","last":"Bar","address":"","ip":"1.1.1.1","email":"info@example.com","epoch":1582889588.3967564,"words":["property","truth","across"]}
{"index" : {"_index" : "elastic-dm-test", "_id" : "abcde" } }
{"name":"Foo","last":"Bar","address":"","ip":"1.1.1.1","email":"info@example.com","epoch":1582889588.3967564,"words":["property","truth","across"]}
'
