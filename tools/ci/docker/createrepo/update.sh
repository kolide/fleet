#!/bin/bash
# create new snapshot folder
mkdir -p /newrepo

# copy current conents to new snapshot
cp -R /repo/* /newrepo

# run createrepo to re-generate metadata
createrepo --update /newrepo

# replace original repodata with snapshot
rm -rf /repo/repodata/*
mv /newrepo/repodata/* /repo/repodata
