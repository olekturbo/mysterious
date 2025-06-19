FROM alpine:latest
COPY mysterious /mysterious
RUN chmod +x /mysterious
ENTRYPOINT ["/mysterious"]
