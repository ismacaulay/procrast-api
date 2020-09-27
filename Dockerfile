FROM golang:1.14-alpine as builder
RUN mkdir /build
WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build ./cmd/procrast-api && go build ./cmd/admin

FROM alpine
RUN adduser -S -D -H -h /app appuser
USER appuser
COPY --from=builder /build/procrast-api /build/admin /app/
WORKDIR /app
CMD ["./procrast-api"]
