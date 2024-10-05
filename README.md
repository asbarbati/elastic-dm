# ElasticSearch Data Migration
The primary objective of this project is to develop a Golang-based tool as an alternative to ElasticDump, aiming to achieve superior performance and incorporate some useful functionalities.
Note: **This tool currently its not ready for the production. Use this tool at your own risk.**

## How to use it?
1. To compile the binary use `go build -o elastic-dm` or download the latest from release page.
2. Copy the file `config.sample.json` into `config.json`.
3. Customize the application settings by editing the config file. Refer to the section below for a complete list of variables and their descriptions.
4. Run it using `./elastic.dm -config config.json`.
5. Watch the log until completed.

## Configuration

|   Name                    |    Type   |   Default  |      Description         |
|   :---:                   |   :---:   |   :---:    |      :---:               |
|   es_src                  |   object  |            | The main object where store all the information about the `SOURCE` elasticsearch |
|   es_src.host             |   string  |            | Elasticsearch Host to use, that needs to have schema on it    |
|   es_src.disabletlsverify |    bool   |   false    | Disable TLS Certification verification    |
|   es_src.user             |   string  |            | User to login into Elasticsearch                     |
|   es_src.passwd           |   string  |            | Password to login into Elasticsearch                 |
|   es_src.indices          |   string  |            | Specify the indices to be processed by the tool. Use commas to separate multiple indices. |
|   es_dst                  |   object  |            | Same objects as `es_src` but for the `DEST` target |
|   mode                    |   string  |   None     | The action to do |
|   scrollmultiplier        |   int     |    10      | The scroll multiplier used in bulk requests  |
|   workers                 |   int     |    10      | The number of concurrent worker goroutines to use for bulk requests    |

## Logging
By default, the tool logs all activity in JSON format to the file elastic-dm.json.log. This facilitates integration with other tools.
You can improve the debugging level adding `-debug` flag.

## Run tests
You can run test using this command below:

```
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html
```
Checkout the file `coverage.html`.