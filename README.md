# Kolide

[![Build Status](https://travis-ci.com/kolide/kolide-ose.svg?token=MvaZkzWisgsA98PZfNC7&branch=master)](https://travis-ci.com/kolide/kolide-ose)

## Building

To build the code ensure you have `node` and `npm` installed run the
following from the root of the repository:

```
npm install
make
```

This will produce a binary called `kolide` in the root of the repo.

## Testing

To run the application's tests, run the following from the root of the
repository:

```
go test
```

Or if you using the Docker development environment run:

```
docker-compose app exec go test
```

## Development Environment

To setup a working local development environment run perform the following tasks:

1. Install the following dependencies:
  * [Docker & docker-compose](https://www.docker.com/products/overview#/install_the_platform)
  * [go 1.6.x](https://golang.org/dl/)
  * [nodejs 0.6.x](https://nodejs.org/en/download/current/) (and npm)
  * A GNU compatible version of `make`

1. Start up all external servers with `docker-compose up`

1. In the root of the repository run:

```
npm install
make
./kolide prepare-db
make serve
```

By default, the last command will run the development proxy on
`http://localhost:8081` which allows you to make live changes to the code and
have them hot-reload.


## Docker Deployment
This repository comes with a simple Dockerfile. You can use this to easily
deploy Kolide in any infrastructure context that can consume a docker image
(heroku, kubernetes, rancher, etc).

To build the image locally, run:

```
docker build --rm -t kolide .
```

To run the image locally, simply run:

```
docker run -t -p 8080:8080 kolide
```
