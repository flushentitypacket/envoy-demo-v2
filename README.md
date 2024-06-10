# Envoy load balancer demo

Demo project built on docker compose, featuring:

- grpc client as load generator (Go)
- Envoy
- grpc server (Go)
- dogstatsd

## Setup

Must have Docker installed.

Copy and modify `.env.example` file to enable submitting metrics to Datadog.

```sh
cp .env.example .env
```

## Usage

```sh
docker compose up --force-recreate --build
```

## Development

### Compile proto

```sh
protoc --go_out=lib/proto --go-grpc_out=lib/proto lib/proto/dummy/dummy.proto
```
