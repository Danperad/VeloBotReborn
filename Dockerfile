FROM golang:1.23-alpine AS build-stage

WORKDIR /app

COPY *.go ./
COPY controllers ./controllers
COPY dbModels ./dbModels
COPY repositories ./repositories
COPY utils ./utils

COPY go.mod go.sum ./
RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build -o /velo-bot

FROM alpine:latest

COPY --from=golang:1.23-alpine /usr/local/go/ /usr/local/go/
ENV PATH="/usr/local/go/bin:${PATH}"

COPY --from=build-stage /velo-bot /velo-bot

ENTRYPOINT ["/velo-bot"]
