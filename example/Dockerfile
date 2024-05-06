FROM golang:latest

RUN apk add --no-cache git

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64 \

WORKDIR /build

COPY go.sum .
COPY go.mod .

RUN go mod download

COPY ../zeus .

RUN go build -o main .

FROM scratch

COPY --from=builder /build/main /

COPY config /config

USER nobody

ENTRYPOINT ["/main"]
