FROM node:20-alpine as webbuilder

WORKDIR /app

COPY . .

RUN npm ci && npm install -g @go-task/cli && task build:css

FROM golang:alpine as builder

RUN apk add -U --no-cache ca-certificates

WORKDIR /app

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o web cmd/player/main.go

FROM scratch

WORKDIR /app

COPY --from=webbuilder /app/public /app/public
COPY --from=webbuilder /app/templates /app/templates
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /app/web /usr/bin/

ENV GIN_MODE=release

EXPOSE 8081

ENTRYPOINT ["web"]