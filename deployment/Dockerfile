# Builder stage
FROM golang:1.23.3-alpine AS builder
RUN apk add --no-cache bash make
WORKDIR /pagu
COPY .. .
RUN make build

FROM alpine:3.14 AS pagu
WORKDIR /bin
COPY --from=builder /pagu/build/pagu-cli .
COPY --from=builder /pagu/build/pagu-discord .
COPY --from=builder /pagu/build/pagu-telegram .
