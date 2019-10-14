FROM golang:alpine as builder
RUN mkdir /build 
ADD . /build/
WORKDIR /build 
RUN go build -o main .
FROM alpine
RUN adduser -S -D -H -h /app appuser
USER appuser
COPY --from=builder /build/main /app/
ADD ./config/user-config.yaml /app/config/user-config.yaml
WORKDIR /app
CMD ["./main", "-yc", "./config/user-config.yaml"]