FROM ubuntu

RUN apt-get update -y && apt-get install -y aptly
COPY update.sh /update.sh
CMD ["/update.sh"]
