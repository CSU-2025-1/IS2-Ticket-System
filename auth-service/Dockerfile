ARG GO_VERSION=1.24

FROM golang:${GO_VERSION}-alpine3.21 AS builder

RUN mkdir /app
WORKDIR /app

COPY ./go.mod ./go.sum ./
RUN go mod download -x

COPY . .

RUN go build -a -ldflags='-s' -trimpath -o /app/bin/app /app/cmd/app/main.go

FROM scratch

COPY --from=builder /app/bin/app ./
COPY ./migrations ./migrations

EXPOSE 8080

ENTRYPOINT ["./app"]
