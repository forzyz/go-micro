# Use a recent Go version and Alpine base image
FROM golang:1.23-alpine AS build

# Install necessary packages
RUN apk --no-cache add gcc g++ make ca-certificates

# Set the working directory
WORKDIR /go/src/github.com/forzyz/go-micro

# Copy the go.mod and go.sum files
COPY go.mod go.sum ./

# Copy the vendor directory
COPY vendor vendor

# Copy the source code
COPY account account
COPY catalog catalog
COPY order order
COPY graphql graphql

# Build the Go application
RUN GO111MODULE=on go build -mod vendor -o /go/bin/app ./graphql

# Use a minimal Alpine image for the final stage
FROM alpine:3.18

# Set the working directory for the final image
WORKDIR /usr/bin

# Copy the compiled binary from the build stage
COPY --from=build /go/bin/app .

# Expose the application port
EXPOSE 8080

# Set the command to run the application
CMD ["app"]
