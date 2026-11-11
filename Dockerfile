# build
FROM golang:1.26-bookworm AS builder
LABEL authors="support@etzba.com, Nadav Ben Mazia"
# copy go mod and sum for better caching
COPY . /build
WORKDIR /build
# Go mod download and verify dependencies
RUN go mod download
RUN go mod verify
# Build the app
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags='-w -s -extldflags "-static"' \
    -a -installsuffix cgo \
    -o gosqlite main.go
#RUN apt update && apt install sqlite3

#FROM gcr.io/distroless/static-debian13
FROM alpine:latest
RUN apk update && apk add sqlite
# Copy the app from builder
COPY --from=builder /build/gosqlite /gosqlite
#COPY --from=builder /usr/bin/sqlite3 /bin/sqlite
WORKDIR /
CMD ["./gosqlite"]
