FROM ubuntu
RUN apt-get update -y && apt-get install -y \
    curl git apt-transport-https build-essential && \
    curl -sS https://dl.yarnpkg.com/debian/pubkey.gpg | apt-key add - && \
    echo "deb https://dl.yarnpkg.com/debian/ stable main" | tee /etc/apt/sources.list.d/yarn.list && \
    apt-get update -y && apt-get install -y yarn

RUN curl -L https://storage.googleapis.com/golang/go1.8.linux-amd64.tar.gz -o go1.8.linux-amd64.tar.gz && \
    tar -C /usr/local -xzf go1.8.linux-amd64.tar.gz && \
    rm go1.8.linux-amd64.tar.gz

ENV GOPATH /go
ENV PATH $GOPATH/bin:/usr/local/go/bin:$PATH

RUN mkdir -p "$GOPATH/src/github.com/kolide/kolide" "$GOPATH/bin" && chmod -R 777 "$GOPATH"
WORKDIR $GOPATH/src/github.com/kolide/kolide

