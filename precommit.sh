#!/bin/sh

TMP=`mktemp`

for PKG in `go list github.com/jollheef/henhouse/... | tr '\n' ' '`; do
    figlet -f big `echo ${PKG} | sed 's|[.a-Z]*/[.a-Z]*/||'`
    echo '---------------' LINT '---------------'
    golint ${PKG} | tee ${TMP} \
        && cat ${TMP} | wc -l | grep 0  >/dev/null && echo "ok" || exit 1
    echo
    echo '---------------' TEST '---------------'
    go test -v -covermode=count -coverprofile=coverage.out ${PKG} || exit 1
    echo
    echo
done

rm ${TMP}
