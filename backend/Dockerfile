FROM golang:alpine as builder
ADD . /go/src/webserver
WORKDIR /go/src/webserver
RUN CGO_ENABLED=0 GOOS=linux go build -v ./cmd/apiserver
RUN apk update && apk add --no-cache git ca-certificates && update-ca-certificates

FROM scratch
COPY --from=builder /go/src/webserver/apiserver /go/bin/apiserver
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /go/src/webserver/configs/server.toml /go/bin/configs/server.toml
EXPOSE 8080
ENTRYPOINT ["/go/bin/apiserver"]

