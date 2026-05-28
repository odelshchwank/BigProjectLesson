FROM golang:1.26.3-bookworm AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /app/todoapp /app/cmd/todoapp

FROM alpine:3.23
WORKDIR /app
COPY --from=builder /app/todoapp /app/todoapp
CMD ["/app/todoapp"]