
FROM golang:1.24-alpine AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy dependency tracking files first to optimize layer caching
COPY go.mod ./
RUN go mod download


COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /ascii-art-web-server .

FROM alpine:3.20

LABEL maintainer="Eddy-Odero" \
      project="dockerize" \
      version="1.0" \
      description="Dockerized Go ASCII Art Web Application"

WORKDIR /

COPY --from=builder /ascii-art-web-server /ascii-art-web-server
COPY --from=builder /app/template /template
COPY --from=builder /app/banners /banners
COPY --from=builder /app/ascii /ascii

EXPOSE 8080

# Run the binary
CMD ["/ascii-art-web-server"]
