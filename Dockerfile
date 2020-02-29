FROM golang:alpine as builder
WORKDIR $GOPATH/src/mypackage/myapp/
COPY . .
#build the binary
RUN CGO_ENABLED=0 go build -tags netgo -a -v ./cmd/apiserver
RUN mkdir /go/bin/configs
COPY configs/server.toml /go/bin/configs/server.toml
COPY apiserver /go/bin/apiserver

FROM scratch
#COPY --from=builder /app/apiserver /app/apiserver
COPY --from=builder /go/bin/apiserver /go/bin/apiserver
COPY --from=builder /go/bin/configs/server.toml /go/bin/configs/server.toml
ENTRYPOINT ["/go/bin/apiserver"]

