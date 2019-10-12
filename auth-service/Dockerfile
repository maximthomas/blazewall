FROM golang:alpine as builder
RUN mkdir /build 
ADD . /build/
WORKDIR /build 
RUN go build -o main .
FROM alpine
RUN adduser -S -D -H -h /app appuser
USER appuser
COPY --from=builder /build/main /app/
ADD ./static/ /app/static
ADD ./templates/ /app/templates
ADD ./config/auth-config.yaml /app/config/auth-config.yaml
WORKDIR /app
CMD ["./main", "-yc", "./config/auth-config.yaml"]