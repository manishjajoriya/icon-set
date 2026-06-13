FROM golang:1.26.4-alpine3.23 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o server ./cmd

# Run stage

FROM alpine:3.23.3

RUN adduser -D -g '' appuser

WORKDIR /app

COPY --from=builder /app/server .
COPY /internal/assets/ ./internal/assets/
COPY allowed_icon_pack.json .

USER appuser

CMD ["./server"]