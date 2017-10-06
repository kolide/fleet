FROM golang:1.9-alpine
MAINTAINER Kolide Developers <engineering@kolide.co>

RUN apk add --update \
        git \
        nodejs \
	bash \
	alpine-sdk \
	musl-dev \
	openssl \
	openssh-client

RUN go get -u -v github.com/Masterminds/glide

RUN echo -e 'http://dl-cdn.alpinelinux.org/alpine/edge/main\nhttp://dl-cdn.alpinelinux.org/alpine/edge/community\nhttp://dl-cdn.alpinelinux.org/alpine/edge/testing' > /etc/apk/repositories && apk add --no-cache yarn
RUN apk add --no-cache nodejs

# git must be explicitly instructed to follow redirects for compatibility with
# gopkg.in
RUN git config --global http.https://gopkg.in.followRedirects true

COPY rootfs /

VOLUME   /go/src/github.com/kolide/fleet
WORKDIR  /go/src/github.com/kolide/fleet

EXPOSE 8080

ENTRYPOINT ["/builder.sh"]
