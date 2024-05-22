# Auth Gateway

![License](https://img.shields.io/github/license/voiplens/auth-gateway.svg?color=blue)
[![Go Report Card](https://goreportcard.com/badge/github.com/voiplens/auth-gateway)](https://goreportcard.com/report/github.com/voiplens/auth-gateway)
![GitHub release](https://img.shields.io/github/v/release/voiplens/auth-gateway)

Auth Gateway is a microservice which strives to help you administrating and operating your [Mimir](https://github.com/grafana/mimir) & [Loki](https://github.com/grafana/loki) Cluster in multi tenant environments.

## Features

- [x] Authentication of Prometheus & Grafana instances with JSON Web Tokens
- [x] Prometheus & Jager instrumentation, compatible with the rest of the Mimir microservices

#### Authentication Feature

If you run Mimir for multiple tenants you need to identify your tenants every time they send metrics or query them. This is needed to ensure that metrics can be ingested and queried separately from each other. For this purpose the Mimir microservices require you to pass a Header called `X-Scope-OrgID`. Unfortunately the Prometheus Remote write API has no config option to send headers and for Grafana you must provision a datasource to do so. Therefore the Mimir k8s manifests suggest deploying an NGINX server inside of each tenant which acts as reverse proxy. It's sole purpose is proxying the traffic and setting the `X-Scope-OrgID` header for your tenant.

We try to solve this problem by adding a Gateway which can be considered the entry point for all requests towards Mimir (see [Architecture](#architecture)). Prometheus and Grafana can both send a self contained JSON Web Token (JWT) along with each request. This JWT carries a claim which is the tenant's identifier. Once this JWT is validated we'll set the required `X-Scope-OrgID` header and pipe the traffic to the upstream Mimir microservices (distributor / query frontend).

# Mimir Auth Gateway

## Architecture

![Mimir Gateway Architecture](./docs/imgs/architecture.png)

## Configuration

| Flag                              | Description                                                                                            | Default        |
| --------------------------------- | ------------------------------------------------------------------------------------------------------ | -------------- |
| `-gateway.distributor.address`    | Upstream HTTP URL for Mimir Distributor                                                                | (empty string) |
| `-gateway.query-frontend.address` | Upstream HTTP URL for Mimir Query Frontend                                                             | (empty string) |
| `-gateway.rules.address`          | Upstream HTTP URL for Mimir ruler                                                                      | (empty string) |
| `-gateway.alertmanager.address`   | Upstream HTTP URL for Mimir Alertmanager                                                               | (empty string) |
| `-gateway.auth.jwt-secret`        | HMAC secret to sign JSON Web Tokens                                                                    | (empty string) |
| `-gateway.auth.jwt-extra-headers` | A comma supported list of additional headers to scan for JSON web tokens presence                      | (empty string) |
| `-gateway.auth.tenant-name`       | The tenant name to use when you want to disable jwt auth, if specified the jwt secret value is ignored | (empty string) |
| `-gateway.auth.tenant-id-claim`   | The name of the tenant ID Claim                                                                        | 'tenant_id'    |

### Expected JWT payload

The expected Bearer token payload can be found here: https://github.com/rewe-digital/cortex-gateway/blob/b74de65d10a93e1ec0d223e92c08d16d59bbf3c4/gateway/tenant.go#L7-L11

- "tenant_id"
- "aud"
- "version" (must be an integer)

The audience and version claim is currently unused, but might be used in the future (e. g. to invalidate tokens).

# Loki Auth Gateway

## Configuration

| Flag                              | Description                                                                                            | Default        |
| --------------------------------- | ------------------------------------------------------------------------------------------------------ | -------------- |
| `-gateway.distributor.address`    | Upstream HTTP URL for Loki Distributor                                                                 | (empty string) |
| `-gateway.query-frontend.address` | Upstream HTTP URL for Loki Query Frontend                                                              | (empty string) |
| `-gateway.rules.address`          | Upstream HTTP URL for Loki Ruler                                                                       | (empty string) |
| `-gateway.querier.address`        | Upstream HTTP URL for Loki Querier                                                                     | (empty string) |
| `-gateway.auth.jwt-secret`        | HMAC secret to sign JSON Web Tokens                                                                    | (empty string) |
| `-gateway.auth.jwt-extra-headers` | A comma supported list of additional headers to scan for JSON web tokens presence                      | (empty string) |
| `-gateway.auth.tenant-name`       | The tenant name to use when you want to disable jwt auth, if specified the jwt secret value is ignored | (empty string) |
| `-gateway.auth.tenant-id-claim`   | The name of the tenant ID Claim                                                                        | 'tenant_id'    |

Forked from: [https://github.com/rewe-digital/cortex-gateway](https://github.com/rewe-digital/cortex-gateway)

# How to run locally

Simply build and run with following commands:

```
make clean all

# Run the mimir auth gateway
cmd/mimir-gateway/mimir-gateway

# Run the loki auth gateway
cmd/loki-gateway/loki-gateway
```

### How to release

Add a new release tag and build and publish the docker image. Commands to build and publish image:

```
docker image build . -f ./Dockerfile --tag $CONTAINER_REGISTRY/cortex-gateway:v0.0.0
docker push $CONTAINER_REGISTRY/cortex-gateway:v0.0.0
```
