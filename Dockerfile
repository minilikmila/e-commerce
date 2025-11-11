FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o /ecommerce ./cmd/server

FROM alpine:3.18

# It's good practice to run as a non-root user
RUN addgroup -S appgroup && adduser -S appuser -G appgroup
USER appuser

WORKDIR /app
# COPY --from=builder /ecommerce /usr/local/bin/ecommerce
# Copy the static binary from the builder stage
COPY --from=builder /ecommerce .

EXPOSE 8080

ENTRYPOINT ["./ecommerce"]

