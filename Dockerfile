FROM alpine:3.4
MAINTAINER Victor Vrantchan <victor@kolide.co> (@groob)

COPY ./build/kolide /kolide

CMD ["/kolide", "serve"]
