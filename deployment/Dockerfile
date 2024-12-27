# Builder stage
FROM golang:1.23.3-alpine AS builder
RUN apk add --no-cache bash make
WORKDIR /pagu
COPY .. .
RUN make build-discord build-telegram

# Discord bot
FROM alpine:3.14 AS discord
WORKDIR /bin
COPY --from=builder /pagu/build/pagu-discord .
CMD ["./pagu-discord", "-c", "/pagu/config.yml", "run"]

# Telegram bot
FROM alpine:3.14 AS telegram
WORKDIR /bin
COPY --from=builder /pagu/build/pagu-telegram .
CMD ["./pagu-telegram", "-c", "/pagu/config.yml", "run"]
