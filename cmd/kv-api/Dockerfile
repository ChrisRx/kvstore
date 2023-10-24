FROM golang:1.21-alpine as builder
WORKDIR /app
COPY . .
RUN apk add git make
RUN make kv-api

FROM scratch
COPY --from=builder /app/bin/kv-api /kv-api
COPY --from=alpine:latest /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
ENTRYPOINT ["/kv-api"]
