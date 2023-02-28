FROM golang:1.19-alpine as builder

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY *.go ./

RUN go build -o /data-sim

# build running image
FROM alpine:3.16

COPY --from=builder /data-sim /usr/local/bin/

WORKDIR /usr/local/bin/

ENTRYPOINT ["data-sim"]
