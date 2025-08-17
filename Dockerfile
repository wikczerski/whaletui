FROM golang:1.25-alpine AS builder

RUN apk add --no-cache git ca-certificates

WORKDIR /app
COPY go.mod go.sum ./

RUN go mod download

COPY . .

ARG VERSION=dev
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags="-s -w -X main.Version=${VERSION}" -o d5r .

FROM alpine:latest

RUN apk --no-cache add ca-certificates
RUN addgroup -g 1001 -S d5r && \
    adduser -u 1001 -S d5r -G d5r

WORKDIR /app
COPY --from=builder /app/d5r .

RUN chown d5r:d5r /app/d5r
USER d5r

ENTRYPOINT ["./d5r"]
