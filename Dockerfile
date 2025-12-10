# Build stage
FROM golang:1.24 AS builder

WORKDIR /app

# Copy only go.mod (you have no go.sum)
COPY go.mod ./
RUN go mod download

# Copy the rest of the source
COPY . .

RUN go build -o blnkGoProxy .

# Runtime image
FROM debian:bookworm-slim

WORKDIR /app

COPY --from=builder /app/blnkGoProxy .

ENV TARGET_HOST=blnk_server
ENV TARGET_PORT=5001
ENV PROXY_PORT=5000

EXPOSE 5000

CMD ["./blnkGoProxy"]
