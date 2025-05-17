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

## tcp-echo-server

[TCO-Echo-Server](tcp-echo-erver) runs the most basic TCP server responding with the query.

```sh
go run .
```

In another terminal:

```sh
echo "hello" | nc localhost 8080
```
