FROM golang:1.14-alpine AS swetasksrc

RUN set -ex; \
    apk update; \
    apk add --no-cache git

WORKDIR /go/src/github.com/stormdynamics/swerpctask/
COPY . ./

RUN set -ex; \
    go mod download; \
    GOOS=linux go build -o ./swerpctask;

RUN ./swerpctask migrate up

FROM alpine:latest

RUN set -ex; \
    apk update; \
    apk add --no-cache \
    ca-certificates

WORKDIR /opt/

COPY --from=swetasksrc /go/src/github.com/stormdynamics/swerpctask/swerpctask .

CMD /opt/swerpctask