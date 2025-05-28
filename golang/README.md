# Golang

## logging

[Logging](logging) demonstrates how to use a global logger

## CLI

[CLI](cli) contains useful snippets for creating CLIs

## gitlab-http

[glab-api-http](glab-api-http) demonstrates Gitlab API calls

## manifest-splitter

[Manifest-splitter](manifest-splitter) splits a file with multiple YAML documents consisting of Kubernetes manifest into one file for each `kind`

## postgres-connection

[Postgres-connection](postgres-connection) to verify a working connection to a Postgres database

## echo-server

[Echo-Server](echo-server) runs a http and TCP echo server

```sh
go run .
```

In another terminal:

```sh
echo "hello" | nc localhost 8080
curl -X POST -d "foo" localhost:8081
```
