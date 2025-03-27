FROM golang:1.24 AS builder

RUN apt update && apt install git ca-certificates gcc -y && update-ca-certificates

ENV USER=appuser
ENV UID=10001

RUN adduser \
    --disabled-password \
    --gecos "" \
    --home "/nonexistent" \
    --shell "/sbin/nologin" \
    --no-create-home \
    --uid "${UID}" \
    "${USER}"

WORKDIR /src
COPY . .

RUN go mod download
RUN go mod verify

RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -tags netgo -o /bin/app --ldflags "-extldflags -static" ./cmd/app/main.go

FROM scratch

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/group /etc/group
COPY --from=builder /bin/app /bin/app

USER appuser:appuser
ENTRYPOINT ["/bin/app"]