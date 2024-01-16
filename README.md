# kvstore

[![PkgGoDev](https://pkg.go.dev/badge/github.com/ChrisRx/kvstore)](https://pkg.go.dev/github.com/ChrisRx/kvstore)

kvstore is a Key-Value store composed of two services. The services are both written in Go and integrated using [gRPC](https://grpc.io/).

## Table of Contents

- [Overview](#overview)
- [Getting started](#getting-started)
  - [Using Docker Compose](#using-docker-compose)
  - [Using Nix](#using-nix)
    - [Install Nix](#installing-nix)
    - [Build the Nix development environment](#build-the-nix-development-environment)
  - [Building and running manually](#building-and-running-manually)
  - [Testing with curl](#testing-with-curl)
- [Development](#development)
  - [New storage backends](#new-storage-backends)
  - [Generating protobuf](#generating-protobuf)
  - [Running integration tests](#running-integration-tests)
  - [TLS authentication](#tls-authentication)

## Overview

The two services are provided as static Go binaries, which take configuration via environment variables and command-line arguments/flags. The README.md for each binary details more about configuration:

| Command | Description |
| --- | --- |
| [kv-api](cmd/kv-api/README.md) | A REST API for kvstore, intended to be public-facing component |
| [kv-grpc-server](cmd/kv-grpc-server/README.md) | A gRPC service providing access to the backend KV storage |

Each binary is also built in a multi-stage container that produces [scratch](https://hub.docker.com/_/scratch/) images for the resulting binaries. This ensures the images are small and reduces the surface area for exploitation by including only what is needed.


## Getting started

### Using Docker Compose

[Docker Compose](https://docs.docker.com/compose/) allows for defining and running multi-container applications. The provided [`compose.yaml`](compose.yaml) file is setup to build and run the two services and only requires running `docker compose up`:


```shell
❯ docker compose up
[+] Building 0.0s (0/0)                         docker:default
[+] Running 2/0
 ✔ Container kvstore-kv-api-1          Created            0.0s
 ✔ Container kvstore-kv-grpc-server-1  Created            0.0s
Attaching to kvstore-kv-api-1, kvstore-kv-grpc-server-1
kvstore-kv-grpc-server-1  | Running kv-grpc-server on :9090 ...
kvstore-kv-api-1          |
kvstore-kv-api-1          |    ____    __
kvstore-kv-api-1          |   / __/___/ /  ___
kvstore-kv-api-1          |  / _// __/ _ \/ _ \
kvstore-kv-api-1          | /___/\__/_//_/\___/ v4.11.2
kvstore-kv-api-1          | High performance, minimalist Go web framework
kvstore-kv-api-1          | https://echo.labstack.com
kvstore-kv-api-1          | ____________________________________O/_______
kvstore-kv-api-1          |                                     O\
kvstore-kv-api-1          | ⇨ http server started on [::]:8080
```

When changes are made to the source code, the images can be rebuilt using `docker compose build` and then re-running `docker compose up`

With the two containers now running, [test with curl](#testing-with-curl) and ensure everything is working.


### Using Nix

#### Installing Nix

```shell
curl --proto '=https' --tlsv1.2 -sSf -L https://install.determinate.systems/nix | sh -s -- install
```

#### Build the Nix development environment

Using `nix develop` will use the `flake.nix` in this repo to create a shell with the packages needed to build this project.

```shell
kvstore on  main via 🐹
❯ nix develop -c $SHELL
...

kvstore on  main via 🐹 via ❄️  impure (nix-shell-env)
❯
```

It will build the environment as defined in `flake.nix` and create a new shell. The `-c $SHELL` flag is used if a shell other than bash should be used.

Once the environment is setup, test that it is using the dependencies defined in `flake.nix` by running `make generate` and it should use the Go and protobuf versions available in nixpkgs.


### Building and running manually

Running `make all` will build static Go binaries for each of the two components and place them in the `bin/` folder:

```shell
Permissions Size User  Date Modified Name
drwxr-xr-x     - chris 25 Oct 10:07  bin/
.rwxr-xr-x   18M chris 25 Oct 11:07  ├── kv-api*
.rwxr-xr-x   15M chris 25 Oct 10:07  └── kv-grpc-server*
```

Next start the gRPC server:


```shell
❯ bin/kv-grpc-server
Running kv-grpc-server on :9090 ...
```

CLI arguments/flags can be used to change the listener address, or specify TLS server authentication, although this will default to `:9090` and without auth. The name of the [boltdb](https://github.com/etcd-io/bbolt) file can also be specified, but defaults to simply `data.db`.

In another terminal, start the REST API:

```shell
❯ bin/kv-api --insecure

   ____    __
  / __/___/ /  ___
 / _// __/ _ \/ _ \
/___/\__/_//_/\___/ v4.11.2
High performance, minimalist Go web framework
https://echo.labstack.com
____________________________________O/_______
                                    O\
⇨ http server started on [::]:8080
```

The `--insecure` flag must be present when not using TLS authentication. This is an important safety measure to protect against accidental production configuration that excludes authentication.

Finally, make sure the services are running correctly using [curl](testing-with-curl):


### Testing with curl

Basic testing can be performed using curl:

```shell
❯ curl -XPOST -H 'Content-Type: application/json' localhost:8080/testkey --data-binary '{"value": "somedata"}'
{"status":"OK"}
```

The value that key `testkey` should now be set:

```shell
❯ curl localhost:8080/testkey
{"key":"testkey","value":"somedata"}
```


## Development


### New storage backends

The KV gRPC service is defined by [internal/kvpb/kv.proto](internal/kvpb/kv.proto), which generates the `KVServer` interface. This interface must be fulfilled by all implementations that are to be registered with the gRPC server that hosts the KV service. The default implementation using boltdb can be found in the [boltkv](internal/boltkv) pacakge and can be used as an example to create new packages that fulfill the `KVServer` interface. This can enable the underlying storage of the KV service to be changed to use other databases/services for backend storage.


### Generating protobuf

By default, gRPC using [protocol buffers](https://protobuf.dev/) as an Interface Definition Language (IDL) and for message serialization, and kvstore uses the official Golang implementation of [protobuf](https://github.com/golang/protobuf). This means that when the proto definition file(s) change, they must be regenerated.

This requires having the protocol buffer compiler (`protoc`) installed, and the Go-specific plugins that are needed. The Go plugins can be installed into the current Go workspace by running `make generate-install` (this only needs to be done once to ensure the plugins are installed). The updated files can then be generated by simply running `make generate`.


### Running integration tests

Some basic integration/end-to-end tests can be found in [test/](test/). After all of the services are [running using Docker Compose](#using-docker-compose):

```shell
❯ go test -v ./test/
=== RUN   TestAPISetDeleteValue
--- PASS: TestAPISetDeleteValue (0.02s)
PASS
ok      github.com/ChrisRx/kvstore/test (cached)
```

### TLS authentication

Any gRPC-based service deployed to production should use TLS authentication. This can be accomplished by specifying the appropriate command-line flags for the kv-api and kv-grpc-server commands, looking something like:

```shell
❯ kv-api --cert-file cert.pem --key-file key.pem --ca-file ca.pem
```
