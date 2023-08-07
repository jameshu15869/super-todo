# syntax=docker/dockerfile:1

# Build from source
FROM golang:1.20 AS build-stage

WORKDIR /go/src

COPY go.mod go.sum ./
RUN go mod download

COPY . ./

RUN CGO_ENABLED=0 GOOS=linux go build -o /service-todo

FROM alpine:latest AS build-release-stage

WORKDIR /

COPY --from=build-stage /service-todo /service-todo

EXPOSE 4002

CMD ["/service-todo"]
