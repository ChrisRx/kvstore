# kv-grpc-server

kv-grpc-server is a Key-Value gRPC service.

```shell
❯ kv-grpc-server -h
Run KV service gRPC server

Usage:
  kv-grpc-server [flags]

Flags:
      --addr string        Address for gRPC server (default ":9090")
      --ca-file string     TLS auth CA cert file
      --cert-file string   TLS auth cert file
      --db-file string     Name of boltdb file (default "data.db")
  -h, --help               help for kv-grpc-server
      --key-file string    TLS auth key file
```
