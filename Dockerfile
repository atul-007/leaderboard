# Dockerfile
FROM golang:alpine as builder
RUN apk update && apk add --no-cache git
WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o main .

FROM alpine
RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*
WORKDIR /root/
COPY --from=builder /app/main .
CMD ["./main"]
