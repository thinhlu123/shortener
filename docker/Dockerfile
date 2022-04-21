FROM golang:1.16

WORKDIR /app
COPY ./ /app
RUN go mod tidy
RUN go mod vendor
RUN go run main.go