FROM golang:1.25-alpine AS builder

WORKDIR /app
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o inference ./cmd/inference

FROM scratch

COPY --from=builder /app/inference /inference

EXPOSE 8080

ENTRYPOINT ["/inference"]
