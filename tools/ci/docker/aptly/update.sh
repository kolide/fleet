#!/bin/bash

snapshot="$(date +%s)"
aptly repo add kolide /deb
aptly snapshot create "${snapshot}" from repo kolide
aptly publish drop jessie
aptly publish -gpg-key="EBA2F3F6" --distribution="jessie" snapshot "${snapshot}"
