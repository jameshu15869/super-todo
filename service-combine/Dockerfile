# syntax=docker/dockerfile:1

# Build from source
FROM golang:1.20 AS build-stage

WORKDIR /go/src

COPY go.mod go.sum ./
RUN go mod download

COPY . ./

RUN CGO_ENABLED=0 GOOS=linux go build -o /service-combine

FROM alpine:latest AS build-release-stage

WORKDIR /

COPY --from=build-stage /service-combine /service-combine

EXPOSE 4003

CMD ["/service-combine"]
