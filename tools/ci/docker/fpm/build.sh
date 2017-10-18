#!/bin/bash

fpm -s dir -t deb --deb-no-default-config-files -n "fleet" -v ${KOLIDE_VERSION} /pkgroot/usr/=/usr
fpm -s dir -t rpm -n "fleet" -v ${KOLIDE_VERSION} /pkgroot/usr/=/usr
mv fleet* /out

# sign packages
rpmVersion="$(echo ${KOLIDE_VERSION}|sed 's/-/_/g')"
rpm --addsign "/out/fleet-${rpmVersion}-1.x86_64.rpm"
debsigs --sign=origin -k ED07465109A40CBAB43B8CA2AE7B6676E205A2C1 "/out/fleet_${KOLIDE_VERSION}_amd64.deb"
