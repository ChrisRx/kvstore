FROM golang:1.21-alpine as builder
WORKDIR /app
COPY . .
RUN apk add git make
RUN make kv-grpc-server

FROM scratch
COPY --from=builder /app/bin/kv-grpc-server /kv-grpc-server
COPY --from=alpine:latest /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
ENTRYPOINT ["/kv-grpc-server"]
