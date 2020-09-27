FROM golang:1.14-alpine as builder
RUN mkdir /build
COPY . /build/
WORKDIR /build
RUN go build ./cmd/procrast-api

FROM alpine
RUN adduser -S -D -H -h /app appuser
USER appuser
COPY --from=builder /build/procrast-api /app/
WORKDIR /app
CMD ["./procrast-api"]
