FROM golang:1.23-alpine AS build-stage

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY **.go ./

RUN go build -o /velo-bot

FROM alpine:latest

COPY --from=golang:1.23-alpine /usr/local/go/ /usr/local/go/
ENV PATH="/usr/local/go/bin:${PATH}"

COPY --from=build-stage /velo-bot /velo-bot

ENTRYPOINT ["/velo-bot"]
