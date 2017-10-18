#!/bin/bash

snapshot="$(date +%s)"
aptly repo create fleet
aptly repo add fleet /deb
aptly snapshot create "${snapshot}" from repo fleet
aptly publish drop jessie
aptly publish -gpg-key="10E85AFB" --distribution="jessie" snapshot "${snapshot}"
