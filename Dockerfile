# MeowCyber — production image (Linux amd64, CGO + SQLite)
FROM golang:1.24-bookworm AS builder
RUN apt-get update && apt-get install -y --no-install-recommends gcc libc6-dev && rm -rf /var/lib/apt/lists/*
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
ENV CGO_ENABLED=1
RUN go build -o /out/meowcyber cmd/server/main.go

FROM debian:bookworm-slim
RUN apt-get update && apt-get install -y --no-install-recommends ca-certificates && rm -rf /var/lib/apt/lists/*
WORKDIR /app
COPY --from=builder /out/meowcyber .
COPY config.yaml config.yaml
COPY web web
COPY tools tools
COPY roles roles
COPY agents agents
COPY skills skills
ENV MEOWCYBER_DATA_DIR=/data
VOLUME ["/data"]
EXPOSE 8080
CMD ["./meowcyber", "-config", "config.yaml", "--https"]
