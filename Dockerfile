FROM alpine:3.3

USER root

RUN sed -i 's|dl-cdn.alpinelinux.org|mirrors.aliyun.com|g' /etc/apk/repositories
RUN apk add --no-cache ca-certificates curl

COPY bin/linux/mirror-proxy mirror-proxy
RUN chmod u+x mirror-proxy

ENTRYPOINT ["./mirror-proxy"]
CMD []
