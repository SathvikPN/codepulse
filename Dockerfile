# syntax=docker/dockerfile:1

# Stage 1: Build the application
FROM golang:1.18-alpine AS builder

WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /codepulse cmd/codepulse/main.go

# Stage 2: Create the runtime image
FROM scratch

COPY --from=builder /codepulse /codepulse

EXPOSE 8080

ENTRYPOINT ["/codepulse"]