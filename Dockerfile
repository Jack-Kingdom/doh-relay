FROM alpine:latest
ENV DEBUG=''

EXPOSE 8000
WORKDIR /root
COPY bin .
CMD ./doh-relay