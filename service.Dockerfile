FROM golang:1.22-bookworm as builder

WORKDIR /app

COPY go.* ./

RUN go mod download

COPY . ./

RUN go build -v -o server ./cmd/app/

FROM debian:bookworm-slim

RUN set -x && apt-get update && DEBIAN_FRONTEND=noninteractive apt-get install -y \
    ca-certificates && \
    rm -rf /var/lib/apt/lists/*

COPY --from=builder /app/server /app/server

COPY ./config/dev.yml /app/

EXPOSE 8080
EXPOSE 2112

CMD [ "/app/server", "--config", "/app/dev.yml" ]
