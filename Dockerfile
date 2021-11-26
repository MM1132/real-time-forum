# syntax=docker/dockerfile:1

##
## Build
##
FROM golang:1.17 AS build

# Dependancies
WORKDIR /go/src/forum

COPY go.mod ./
COPY go.sum ./
RUN go mod download

# Copy source files
COPY . ./

# Build
RUN go build -o /bin/forum

##
## Deploy
##
FROM gcr.io/distroless/base-debian11
COPY --from=build /bin/forum /forum
COPY server/ /server
EXPOSE 80
LABEL org.opencontainers.image.authors="urmas.rist@gmail.com"
ENTRYPOINT ["/forum"]
