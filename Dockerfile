FROM golang:1.25-alpine AS builder

RUN apk add --no-cache git ca-certificates

WORKDIR /app
COPY go.mod go.sum ./

RUN go mod download

COPY . .

ARG VERSION=dev
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags="-s -w -X main.Version=${VERSION}" -o whaletui .

FROM alpine:latest

RUN apk --no-cache add ca-certificates
RUN addgroup -g 1001 -S whaletui && \
    adduser -u 1001 -S whaletui -G whaletui

WORKDIR /app
COPY --from=builder /app/whaletui .

RUN chown whaletui:whaletui /app/whaletui
USER whaletui

ENTRYPOINT ["./whaletui"]
