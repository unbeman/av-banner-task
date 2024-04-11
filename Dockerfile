FROM golang:1.22.0-alpine
RUN mkdir /app
ADD . /app
WORKDIR /app
RUN go build -o banner-service cmd/main.go