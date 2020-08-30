# solid-broccoli 

Simple Go application that provides HTTP API service

## Public API 

- **/v1/summary/<domain-name>** - returns a count of positions for domain

- **/v1/positions/<domain-name>?orderBy=<field-to-order-by>&page=<page-number>** - returns a bunch of 
positions for domain 

By default, it listens at port 63100.

## Service API

Service API provides standard [pprof](https://golang.org/pkg/net/http/pprof/) endpoints.

By default, it listens at port 63101.

## Build 

Use the following command to build binary:

```bash
make build
```

Or you can use Docker and build image:

```bash 
docker build .
```

## Run

Example of config file: [solid-broccoli.yaml](solid-broccoli.example.yaml)

```bash
./solid-brocoli --config <path-to-yaml-config>
```

## Testing

Use the following command to run acceptance tests (you will need `docker-compose`):

```sh
make acc-tests
```

Use the following command to run only unit-tests and linters:

```sh
make tests
```

Use the following command to run only unit-tests:

```sh
make unittests
```

## Linters

Use the following command to run golangci-lint:

```sh
make golangci-lint
```

Use the following command to run golangci-lint with unit-tests:

```sh
make tests
```
