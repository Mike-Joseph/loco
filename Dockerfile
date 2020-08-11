FROM alpine:3.6
MAINTAINER Thiago Santos <thiago@waltznetworks.com>

COPY loco loco

EXPOSE 8080

ENTRYPOINT ["/loco"]

