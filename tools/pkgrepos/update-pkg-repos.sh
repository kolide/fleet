#!/bin/bash

USER="$(whoami)"
LOCAL_REPO_PATH="/Users/${USER}/kolide_packages"
GPG_PATH="/Users/${USER}/.gnupg"

build_createrepo_container() {
    cd ../ci/docker/createrepo && \
        docker build -t createrepo .
}

build_aptly_container() {
    cd ../ci/docker/aptly && \
        docker build -t aptly .
}

update_yum_repo() {
    # generate new yum repo snapshot
    docker run -it --rm \
        -v "${LOCAL_REPO_PATH}/rpm:/repo" \
        -v "${LOCAL_REPO_PATH}/centos:/repo/repodata" \
        createrepo
    
    # remove artifact from mounting folders to this path
    rm -rf "${LOCAL_REPO_PATH}/rpm/repodata"
    
}

update_apt_repo() {
    docker run -it --rm \
        -v "${LOCAL_REPO_PATH}/deb:/deb" \
        -v "${GPG_PATH}:/root/.gnupg" \
        -v "${LOCAL_REPO_PATH}/aptly:/root/.aptly" \
        -v "${LOCAL_REPO_PATH}/aptly.conf:/root/.aptly.conf" aptly

    # replace "debian" repo with updated snapshot
    rm -rf "${LOCAL_REPO_PATH}/debian" 
    mv "${LOCAL_REPO_PATH}/aptly/public" "${LOCAL_REPO_PATH}/debian" 
}


main() {
    # build_createrepo_container
    build_aptly_container
    # update_yum_repo
    update_apt_repo
}

main
