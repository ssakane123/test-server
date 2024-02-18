FROM golang:1.22.0 AS builder

WORKDIR /test-server

COPY ./main.go ./main.go
COPY ./go.mod ./go.mod
RUN go mod tidy \
    && go build .

FROM rockylinux:9

RUN groupadd -r testserver \
    && useradd -rg testserver testserver -m

USER testserver:testserver

WORKDIR /home/testserver

COPY --from=builder /test-server/test-server /usr/local/bin/