#!/bin/bash

./build.sh

VERSION="$(git tag | sed 's/v//')-$(git log --oneline $(git tag) | wc -l)"
PKGDIR=/tmp/henhouse_${VERSION}

rm -rf ${PKGDIR}

mkdir -p ${PKGDIR}/{DEBIAN,etc/henhouse,usr/bin,lib/systemd/system,var/lib/henhouse,var/www/henhouse}

cp ${GOPATH}/bin/{henhouse,henhousectl} ${PKGDIR}/usr/bin/

cp ./deb/{control,postinst} ${PKGDIR}/DEBIAN/

cp ./henhouse.service ${PKGDIR}/lib/systemd/system/

cp ./config/tasks/bar1.xml ${PKGDIR}/etc/henhouse/example.xml

cp -r ./scoreboard/www/* ${PKGDIR}/var/www/henhouse/

cp ./config/henhouse.toml ${PKGDIR}/etc/

cp -r ./scoreboard/templates ${PKGDIR}/var/lib/henhouse/

sed -i "s/VERSION_PLACEHOLDER/${VERSION}/" ${PKGDIR}/DEBIAN/control

fakeroot dpkg-deb --build ${PKGDIR}

package_cloud push jollheef/henhouse/ubuntu/xenial ${PKGDIR}.deb
