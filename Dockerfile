FROM golang:1.25-alpine3.23 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o _awsd_prompt

FROM alpine:3.23

RUN adduser -D -u 1000 awsd

COPY --from=builder /app/_awsd_prompt /usr/local/bin/_awsd_prompt
COPY scripts/_awsd /usr/local/bin/_awsd
COPY scripts/_awsd_autocomplete /usr/local/bin/_awsd_autocomplete

USER awsd
WORKDIR /home/awsd

ENTRYPOINT ["_awsd_prompt"]
# docker run -it -v ~/.aws:/home/awsd/.aws:ro -v ~/.awsd:/home/awsd/.awsd awsd
