FROM alpine:latest as certs
RUN apk --update add ca-certificates

# build workspace
FROM golang:1.13 as build

WORKDIR /go/src/github.com/whalebone/go-dockerhub-ci
RUN mkdir -p /build

COPY go.mod ./
COPY go.sum ./

COPY *.go ./

# For scratch prod builds
RUN CGO_ENABLED=0 GOOS=linux go build

# prod build
FROM scratch

WORKDIR /go/bin
ENV PATH=/bin
COPY --from=certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=build /go/src/github.com/whalebone/go-dockerhub-ci/go-dockerhub-ci .

ENTRYPOINT [ "./go-dockerhub-ci" ]
