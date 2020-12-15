FROM golang:alpine as build

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

RUN apk --no-cache add ca-certificates

WORKDIR $GOPATH/src/fp-dim-ingress-file-importer/

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -ldflags="-w -s" -o /go/bin/fp-import

FROM scratch AS release

COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build /go/bin/fp-import /
COPY --from=build /go/src/fp-dim-ingress-file-importer/config/ /config/

ENTRYPOINT ["/fp-import"]