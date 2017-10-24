#!/bin/bash
# run createrepo to re-generate metadata
createrepo --update /repo

# sign repo with GPG key
gpg --default-key 10E85AFB --detach-sign --armor /repo/repodata/repomd.xml
