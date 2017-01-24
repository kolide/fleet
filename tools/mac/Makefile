PKGNAME=kolide-enroll
PKGVERSION=1.0.0
PKGID=co.kolide.osquery.enroll

-include config.mk

export KOLIDE_FLAGS
export KOLIDE_TLS_CERTIFICATE


all: clean build
clean:
	rm -rf out/
	rm -rf root/etc/osquery

build: clean
	mkdir -p out
	mkdir -p root/etc/osquery
	echo $(ENROLL_SECRET) > root/etc/osquery/kolide_secret
	echo "$$KOLIDE_TLS_CERTIFICATE" > root/etc/osquery/kolide.crt
	echo "$$KOLIDE_FLAGS" > root/etc/osquery/kolide.flags
	pkgbuild --info PackageInfo --root root --identifier ${PKGID} --version ${PKGVERSION} out/${PKGNAME}-${PKGVERSION}.pkg
