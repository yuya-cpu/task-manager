FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -o server .

FROM alpine:3.21

RUN apk add --no-cache ca-certificates

WORKDIR /app

COPY --from=builder /app/server .
COPY --from=builder /app/web ./web

ENV DB_PATH=/app/data/task-manager.db
ENV SECRET_KEY=change-me-in-production
ENV GIN_MODE=release

RUN mkdir -p /app/data

EXPOSE 8080

CMD ["./server"]
