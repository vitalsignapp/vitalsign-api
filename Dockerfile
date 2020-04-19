FROM golang:1.14.1-alpine3.11 AS build-env
WORKDIR /app
ENV GO111MODULE=on
ENV CGO_ENABLED=0
COPY go.mod .
COPY go.sum .
RUN apk add --no-cache git
RUN go mod download

COPY . /app

RUN go test ./...
RUN GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -tags netgo -ldflags '-w' -o app .

  # final stage
FROM alpine:3.9
WORKDIR /app
RUN apk update && apk add --no-cache git ca-certificates tzdata && update-ca-certificates
COPY --from=build-env /app/app .
COPY --from=build-env /app/configs configs/
RUN mkdir /app/upload

EXPOSE 1323 1324
ENTRYPOINT ["./app"]