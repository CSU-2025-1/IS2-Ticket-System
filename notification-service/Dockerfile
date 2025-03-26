FROM golang:1.24 AS modules

ADD go.mod go.sum /m/
RUN cd /m && go mod download

FROM golang:1.24 AS builder

COPY --from=modules /go/pkg /go/pkg

RUN adduser -u 10001 goless

COPY .  /app

WORKDIR /app

RUN GOOS=linux GOARCH=amd64 go build -ldflags='-s' -a -trimpath -o ./bin/app ./cmd/app/main.go

FROM ubuntu:20.04

COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/group /etc/group
COPY --from=builder /app/bin/app /app/bin/app

RUN chmod +x /app/bin/app

USER goless:goless

# Проверка на исполняемый файл
CMD ["/app/bin/app"]

