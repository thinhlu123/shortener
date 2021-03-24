FROM golang:1.15-alpine AS builder
RUN mkdir /build
ADD . /build/
WORKDIR /build
RUN go build --mod=vendor -o app-exe .


FROM alpine:latest
RUN mkdir /app
WORKDIR /app
COPY --from=builder /build/app-exe .
ARG env
ARG config
ARG version
ENV version=${version}
ENV env=${env}
ENV config=${config}
EXPOSE 8000
CMD ["/app/app-exe"]