FROM alpine:latest

WORKDIR /app_root

COPY app .

EXPOSE 8080

CMD ["./app"]