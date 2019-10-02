FROM golang:alpine as builder
RUN mkdir /build 
ADD . /build/
WORKDIR /build 
RUN go build -o main .
FROM alpine
RUN adduser -S -D -H -h /app appuser
USER appuser
COPY --from=builder /build/main /app/
ADD ./templates/ /app/templates
ADD ./configs/ /app/configs
ADD ./test/ /app/test
WORKDIR /app
CMD ["./main", "-yc", "./configs/auth-config.yaml"]