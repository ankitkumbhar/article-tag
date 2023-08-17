FROM golang:1.19-alpine as builder

RUN mkdir /app

ADD . /app

WORKDIR /app/cmd

RUN go build -o main .

FROM alpine:latest

RUN mkdir /app

WORKDIR /app

COPY --from=builder /app/cmd/main /app

EXPOSE 8080

ENTRYPOINT ./main