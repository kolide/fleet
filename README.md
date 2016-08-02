# Kolide

[![Build Status](https://drone.io/github.com/kolide/kolide-ose/status.png)](https://drone.io/github.com/kolide/kolide-ose/latest)

## Building

To build the code, run the following from the root of the repository:

```
go build -o kolide
```

## Testing

To run the application's tests, run the following from the root of the
repository:

```
go test
```

## Development Environment

To set up a canonical development environment via docker,
run the following from the root of the repository:

```
docker-compose up
```

This requires that you have docker installed. At this point in time,
automatic configuration tools are not included with this project.

If you'd like to shut down the virtual infrastructure created by docker, run
the following from the root of the repository:

```
docker-compose down
```

Once you `docker-compose up` and are running the databases, build the code
and run the following command to create the database tables:

```
kolide prepare-db
```

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
