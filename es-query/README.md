# ES-Query

## Golang

This script performs an Elasticsearch search query and returns prints the first found document

Requires `ELASTICSEARCH_URL` env variable set to the URL of your Elasticsearch. Authentication is not implemented.

## Python

This script performs a slightly complex search query using Elasticsearch' [pagination](https://www.elastic.co/guide/en/elasticsearch/reference/current/paginate-search-results.html) and dumps the result to stdout or a file.
