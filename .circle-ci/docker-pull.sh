#!/bin/bash

mkdir -p ~/docker
if [[ -e ~/docker/kolide_builder.tar ]]; then 
    docker load -i ~/docker/kolide_builder.tar; 
else
    docker pull kolide/kolide-builder:1.8
    docker save kolide/kolide-builder:1.8 > ~/docker/kolide_builder.tar
fi


if [[ -e ~/docker/redis.tar ]]; then 
    docker load -i ~/docker/redis.tar; 
else
    docker pull redis
    docker save redis > ~/docker/redis.tar
fi


if [[ -e ~/docker/mysql.tar ]]; then 
    docker load -i ~/docker/mysql.tar; 
else
    docker pull mysql:5.7
    docker save mysql > ~/docker/mysql.tar
fi
