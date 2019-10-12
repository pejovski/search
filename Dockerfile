# Start from golang v1.13 base image
FROM golang:1.13 as builder

# Add Maintainer Info
LABEL maintainer="Igor Pejovski <igor.pejovski@limango.de>"

# Set the Current Working Directory inside the container
WORKDIR /go/src/github.com/pejovski/search

# Copy everything from the current directory to the PWD(Present Working Directory) inside the container
COPY . .

# Download dependencies
RUN go get -d -v ./...

# Build the Go app
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /go/bin/search .


######## Start a new stage from scratch #######
FROM alpine:3.10

RUN apk --no-cache add ca-certificates

RUN apk add --no-cache tzdata
ENV TZ Europe/Berlin

WORKDIR /root/

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /go/bin/search .

EXPOSE 8207

CMD ["./search"]