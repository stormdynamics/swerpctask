FROM golang:1.14-alpine AS swetasksrc

# установка git для нормальной работы команды go mod (go get)
RUN set -ex; \
    apk update; \
    apk add --no-cache git

WORKDIR /go/src/github.com/stormdynamics/swerpctask/
COPY . ./

RUN set -ex; \
    go mod download; \
    GOOS=linux go build -o ./swerpctask;

# выполнение миграций
RUN ./swerpctask migrate up

FROM alpine:latest

# установка корневых сертификатов
RUN set -ex; \
    apk update; \
    apk add --no-cache \
    ca-certificates

WORKDIR /opt/

# копируем бинарник приложения в /opt
COPY --from=swetasksrc /go/src/github.com/stormdynamics/swerpctask/swerpctask .

# запуск приложения
CMD /opt/swerpctask