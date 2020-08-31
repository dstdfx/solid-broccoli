# solid-broccoli 

Simple Go application that provides HTTP API service

## Public API 

- `/v1/summary/<domain-name>` - returns a count of positions for domain

Example:
```bash
curl -s -X GET "127.0.0.1:63100/v1/summary/fidel.net" | json_pp
{
   "domain" : "fidel.net",
   "positions_count" : 268
}
```

- `/v1/positions/<domain-name>?orderBy=<field-to-order-by>&page=<page-number>` - returns a bunch of
positions for domain 

Example:
```bash
curl -s -X GET "127.0.0.1:63100/v1/positions/fidel.net?page=2" | json_pp
{
   "domain" : "fidel.net",
   "positions" : [
      {
         "cpc" : 1.24,
         "keyword" : "air",
         "results" : 510000000,
         "url" : "https://fidel.net/Theis-enzyms-LAB.html",
         "updated" : "2017-05-23",
         "position" : 41,
         "volume" : 3390000
      },
      {
         "cpc" : 74.85,
         "results" : 2000000000,
         "keyword" : "frame",
         "url" : "https://allodial.fidel.net/",
         "updated" : "2017-05-15",
         "position" : 46,
         "volume" : 4130000
      },
      {
         "cpc" : 74.85,
         "url" : "https://fidel.net/superpure-gnathic-Rowleian.html",
         "keyword" : "frame",
         "results" : 2000000000,
         "position" : 94,
         "updated" : "2017-05-15",
         "volume" : 4130000
      },
      {
         "position" : 74,
         "updated" : "2017-05-14",
         "volume" : 4420000,
         "cpc" : 49.96,
         "results" : 2400000000,
         "keyword" : "card",
         "url" : "https://fidel.net/bestuur"
      },
      {
         "volume" : 4670000,
         "updated" : "2017-05-22",
         "position" : 20,
         "url" : "https://fidel.net/Dyaus",
         "results" : 3100000000,
         "keyword" : "leather",
         "cpc" : 14.6
      },
      {
         "volume" : 5310000,
         "updated" : "2017-05-14",
         "position" : 26,
         "results" : 160000000,
         "keyword" : "flame",
         "url" : "https://Physostomi.fidel.net/",
         "cpc" : 64.43
      },
      {
         "position" : 16,
         "updated" : "2017-05-07",
         "volume" : 5510000,
         "cpc" : 96.91,
         "keyword" : "long",
         "results" : 2200000000,
         "url" : "http://fidel.net/versin/self-govern/tyrantcraft"
      },
      {
         "cpc" : 61.1,
         "keyword" : "public",
         "results" : 23000000,
         "url" : "https://fidel.net/unrefitted",
         "position" : 69,
         "updated" : "2017-05-20",
         "volume" : 5520000
      },
      {
         "updated" : "2017-05-20",
         "position" : 71,
         "volume" : 5520000,
         "cpc" : 61.1,
         "url" : "https://subepiglottal.fidel.net/",
         "keyword" : "public",
         "results" : 23000000
      },
      {
         "url" : "https://fidel.net/unhat",
         "keyword" : "serious",
         "results" : 1100000000,
         "cpc" : 17.28,
         "volume" : 6020000,
         "updated" : "2017-05-08",
         "position" : 14
      }
   ]
}
```

By default, it listens at port 63100.

## Service API

Service API provides standard [pprof](https://golang.org/pkg/net/http/pprof/) endpoints.
By default, it listens at port 63101.

Example of fetching profile over HTTP:
```bash
go tool pprof http://127.0.0.1:63101/debug/pprof/profile
```

You could also visit `http://127.0.0.1:63101/debug/pprof/` in your browser and do some profiling.

It also has `/metrics` endpoint which provides Prometheus metrics.

Example of fetching Prometheus metrics:
```bash
curl -s -X GET 127.0.0.1:63101/metrics
# HELP build_info Build information about the app.
# TYPE build_info gauge
build_info{build_date="",compiler="go1.15",git_commit="",git_tag=""} 1
...
```

## Build 

Use the following command to build binary:

```bash
make build
```

Or you can build Docker image:

```bash 
docker build -t solid-broccoli .
```

## Usage

Example of config file: [solid-broccoli.yaml](solid-broccoli.example.yaml)

Running locally:
```bash
./solid-broccoli --config <path-to-yaml-config>
```

Running in Docker (requires `solid-broccoli` image to be build, see the commands above):
```bash
docker run -p 63101:63101 -p 63100:63100 \
           -v (pwd)/solid-broccoli.yaml:/etc/solid-broccoli/solid-broccoli.yaml \
           solid-brocoli
```

Note, that running the command above you should have `solid-broccoli.yaml` locally.

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
