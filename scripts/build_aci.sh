#!/bin/sh

# Set version of container
VERSION='latest'
if [[ ${#} -eq 1 ]]; then
    VERSION="${1}"
fi

# Generate random machine id for this container
openssl rand -hex 16 > build/machine-id

# Build ACI
acbuild begin
acbuild --debug dependency add quay.io/coreos/alpine-sh
acbuild --debug set-name go-paste
acbuild --debug annotation add version "${VERSION}"
acbuild --debug annotation add author "Lex van Roon <r3boot@r3blog.nl>"
acbuild --debug copy build/machine-id /etc/machine-id
acbuild --debug copy files/resolv.conf /etc/resolv.conf
acbuild --debug copy files/repositories /etc/apk/repositories
acbuild --debug run apk update
acbuild --debug run apk upgrade
acbuild --debug copy build/src/github.com/r3boot/go-paste/build/go-paste /usr/sbin/go-paste
acbuild --debug copy templates/index.html /usr/share/go-paste.html
acbuild --debug set-exec /usr/sbin/go-paste
acbuild --debug write --overwrite ./build/go-paste-${VERSION}-amd64.aci
acbuild end
