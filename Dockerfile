FROM golang:1.16
WORKDIR /
COPY . .
RUN go build --mod=vendor -o app-exe .
CMD ["config=docker", "/app-exe"]