# Builder stage
FROM golang:1.23.6-alpine AS builder
RUN apk add --no-cache bash make
WORKDIR /pagu
COPY .. .
RUN make build

FROM alpine:3.20 AS pagu
WORKDIR /bin
COPY --from=builder /pagu/build/pagu-cli .
COPY --from=builder /pagu/build/pagu-discord .
COPY --from=builder /pagu/build/pagu-telegram .
