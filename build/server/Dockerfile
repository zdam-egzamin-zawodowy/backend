FROM golang:alpine as builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source from the current directory to the Working Directory inside the container
COPY . .
ARG VERSION="0.0.0"
RUN apk --no-cache add musl-dev gcc build-base
RUN go install github.com/99designs/gqlgen@v0.14.0
RUN go generate ./...
RUN go build -ldflags="-X 'main.Version=$VERSION'" -o zdamegzawodowy ./cmd/server

######## Start a new stage from scratch #######
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /app/zdamegzawodowy .

ENV APP_MODE=production
EXPOSE 8080

CMD ./zdamegzawodowy
