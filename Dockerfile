# Create the intermediate builder image.
FROM golang:1.15 as builder

# Docker is copying directory contents so we need to copy them in same directories.
WORKDIR /go/src/github.com/dstdfx/solid-broccoli
COPY cmd cmd
COPY internal internal
COPY vendor vendor
COPY .git .git
COPY data data
# Other files can be copied into the WORKDIR.
COPY ["go*", "./"]

# Build the static application binary.
RUN BUILD_GIT_COMMIT=$(git rev-parse HEAD) \
    BUILD_GIT_TAG=$(git describe --abbrev=0) \
    BUILD_DATE=$(date +%Y%m%d) \
    GO111MODULE=on CGO_ENABLED=1 GOOS=linux \
    go build -mod=vendor -a -installsuffix cgo \
    -ldflags \
    "-X github.com/dstdfx/solid-broccoli/cmd/solid-broccoli/app.buildGitCommit=${BUILD_GIT_COMMIT} \
    -X github.com/dstdfx/solid-broccoli/cmd/solid-broccoli/app.buildGitTag=${BUILD_GIT_TAG} \
    -X github.com/dstdfx/solid-broccoli/cmd/solid-broccoli/app.buildDate=${BUILD_DATE}" \
    -o solid-broccoli ./cmd/solid-broccoli/solidbroccoli.go

# Create the final small image.
FROM alpine:3.12

RUN apk update && apk upgrade && \
    apk add --no-cache \
    ca-certificates wget && \
    rm -rf /var/cache/apk/* && \
    wget -q -O /etc/apk/keys/sgerrand.rsa.pub https://alpine-pkgs.sgerrand.com/sgerrand.rsa.pub && \
    wget https://github.com/sgerrand/alpine-pkg-glibc/releases/download/2.28-r0/glibc-2.28-r0.apk && \
    apk add glibc-2.28-r0.apk

COPY --from=builder /go/src/github.com/dstdfx/solid-broccoli/solid-broccoli /usr/bin/solid-broccoli
COPY --from=builder /go/src/github.com/dstdfx/solid-broccoli/data /data

EXPOSE 63100 63101

ENTRYPOINT ["/usr/bin/solid-broccoli"]
