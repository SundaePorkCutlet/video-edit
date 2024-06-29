# Dockerfile
FROM golang:1.22-alpine AS builder

WORKDIR /bin/app

COPY . .

RUN go mod download

RUN go build -o ./bin/app cmd/app/main.go

FROM alpine:latest

WORKDIR /bin/app

COPY --from=builder /bin/app .

CMD ["./bin/app"]
