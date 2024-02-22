FROM golang:1.22.0 AS builder

WORKDIR /test-server

COPY ./main.go ./main.go
COPY ./go.mod ./go.mod
RUN go mod tidy \
    && go build .

FROM rockylinux:9
ARG CONFIG_FILE="./examples/todo-api.yaml"

RUN groupadd -r testserver \
    && useradd -rg testserver testserver -m

USER testserver:testserver

WORKDIR /home/testserver

COPY ${CONFIG_FILE} ./config.yaml
COPY --from=builder /test-server/test-server /usr/local/bin/

ENTRYPOINT ["test-server", "-config"]
CMD ["config.yaml"]
