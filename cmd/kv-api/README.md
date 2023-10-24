# kv-api

kv-api is a REST api that uses a gRPC client to communicate with a KV service.

```shell
❯ kv-api -h
Run KV REST API server

Usage:
  kv-api [flags]

Flags:
      --addr string        Address for API server (default ":8080")
      --ca-file string     TLS auth CA cert file
      --cert-file string   TLS auth cert file
  -h, --help               help for kv-api
      --insecure           Allow insecure gRPC client connection
      --key-file string    TLS auth key file
      --kv-addr string     Client address for KV gRPC server (default ":9090")
```
