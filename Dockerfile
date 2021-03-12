FROM alpine:latest as certs
RUN apk --update add ca-certificates

# build workspace
FROM golang:1.16 as build

WORKDIR /go/src/github.com/whalebone/go-dockerhub-ci
RUN mkdir -p /build

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

COPY *.go ./

# For scratch prod builds
RUN CGO_ENABLED=0 GOOS=linux go build -o go-dockerhub-ci .

# prod build
FROM scratch

WORKDIR /go/bin
ENV PATH=/bin
COPY --from=certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=build /go/src/github.com/whalebone/go-dockerhub-ci/go-dockerhub-ci .

ENTRYPOINT [ "./go-dockerhub-ci" ]
