FROM golang:1.26-alpine3.23 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o awsd

FROM alpine:3.24

RUN adduser -D -u 1000 awsd

COPY --from=builder /app/awsd /usr/local/bin/awsd

USER awsd
WORKDIR /home/awsd

ENTRYPOINT ["awsd"]
# docker run -it -v ~/.aws:/home/awsd/.aws:ro -v ~/.awsd:/home/awsd/.awsd awsd
