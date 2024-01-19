FROM node:20-alpine as webbuilder

WORKDIR /app

COPY . .

RUN npm ci && npx tailwindcss -i ./assets/css/main.css -o ./cmd/player/assets/css/main.css --minify

FROM golang:alpine as builder

RUN apk add -U --no-cache ca-certificates

WORKDIR /app

COPY --from=webbuilder /app/cmd/player/assets/css/main.css /app/cmd/player/assets/css/main.css

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o web cmd/player/main.go

FROM scratch

WORKDIR /app

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /app/web /usr/bin/

ENV GIN_MODE=release

EXPOSE 8081

ENTRYPOINT ["web"]