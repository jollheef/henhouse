#!/bin/bash

./build.sh

VERSION="$(git tag | sed 's/v//')-$(git log --oneline $(git tag) | wc -l)"
PKGDIR=/tmp/henhouse_${VERSION}

rm -rf ${PKGDIR}

mkdir -p ${PKGDIR}/{DEBIAN,usr/bin}

cp ${GOPATH}/bin/henhouse ${PKGDIR}/usr/bin/
cp ${GOPATH}/bin/henhousectl ${PKGDIR}/usr/bin/

cp ./deb/control ${PKGDIR}/DEBIAN/

sed -i "s/VERSION_PLACEHOLDER/${VERSION}/" ${PKGDIR}/DEBIAN/control

dpkg-deb --build ${PKGDIR}

package_cloud push jollheef/henhouse/ubuntu/xenial ${PKGDIR}.deb
