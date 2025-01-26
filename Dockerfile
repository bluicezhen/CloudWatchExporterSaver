FROM golang:1.21-alpine AS builder

WORKDIR /app
COPY . .
RUN go build -o cloudwatch-exporter-saver

FROM alpine:latest

WORKDIR /app
COPY --from=builder /app/cloudwatch-exporter-saver .

ENV CLOUDWATCH_EXPORTER_URL=""

EXPOSE 9107

CMD ["./cloudwatch-exporter-saver"]
