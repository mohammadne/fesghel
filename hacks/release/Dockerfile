# STEP 1 build executable binary
FROM golang:1.24.2 AS builder
WORKDIR /src
COPY . ./
RUN cd cmd/server && CGO_ENABLED=0 go build -o fesghel && mv fesghel /usr/bin

# STEP 2 build a small image
FROM alpine:3.20
LABEL maintainer="Mohammad Nasr <mohammadne.dev@gmail.com>"
RUN apk add --no-cache bind-tools busybox-extras
COPY --from=builder /usr/bin/fesghel /usr/bin/fesghel
ENV USER=root
ENTRYPOINT ["/usr/bin/fesghel"]
