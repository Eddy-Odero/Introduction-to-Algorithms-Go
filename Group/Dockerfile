# --- Build stage ---
FROM golang:1.22-bookworm AS builder

WORKDIR /app

# go.sum* (with the trailing *) makes this file optional: it won't exist
# until the project has at least one external dependency (e.g. once the
# SQLite driver is added in Phase 3). Docker's COPY treats a wildcard that
# matches zero files as a no-op instead of failing.
COPY go.mod go.sum* ./
RUN go mod download

COPY . .

# CGO is required later for the SQLite driver (mattn/go-sqlite3), so we
# enable it now rather than switching build flags mid-project.
ENV CGO_ENABLED=1
RUN go build -o /forum ./cmd/server

# --- Runtime stage ---
FROM debian:bookworm-slim

# libsqlite3-0 is the runtime library the CGO-based SQLite driver needs.
RUN apt-get update && apt-get install -y --no-install-recommends \
    libsqlite3-0 ca-certificates \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app

COPY --from=builder /forum ./forum
COPY web ./web

EXPOSE 8080

CMD ["./forum"]
