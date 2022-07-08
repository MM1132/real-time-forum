# syntax=docker/dockerfile:1

##
## Build
##
FROM golang:1.17 AS build

# Dependancies
WORKDIR /go/src/real-time-forum

COPY go.mod ./
COPY go.sum ./
RUN go mod download

# Copy source files
COPY . ./

# Build
RUN go build -o /bin/real-time-forum

##
## Deploy
##
FROM gcr.io/distroless/base-debian11
COPY --from=build /bin/real-time-forum /real-time-forum
COPY server/ /server
EXPOSE 80
LABEL org.opencontainers.image.authors="urmas.rist@gmail.com"
ENTRYPOINT ["/real-time-forum"]
