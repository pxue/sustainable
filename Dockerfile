# STAGE 1: BUILD
FROM golang:1.13.0-alpine
ADD . /go/sustainable
WORKDIR /go/sustainable

# GET Signing CERTS
RUN apk update \
        && apk upgrade \
        && apk add --no-cache \
        ca-certificates \
        && update-ca-certificates 2>/dev/null || true

# DEP
RUN apk add --update --no-cache ca-certificates git

# VENDOR
RUN go mod download

# COMPILE
RUN mkdir -p ./bin
RUN GOGC=off CGO_ENABLED=0 GOOS=linux go build -gcflags=-trimpath=${GOPATH} -asmflags=-trimpath=${GOPATH} -a -i -o ./bin/sustainable ./cmd/sustainable/main.go

# STAGE 2: SCRATCH BINARY
FROM scratch
COPY ./db /db
COPY --from=0 /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=0 /go/sustainable/bin/sustainable /bin/sustainable

CMD ["/bin/sustainable"]
