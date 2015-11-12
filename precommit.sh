#!/bin/sh

TMP=`mktemp`

function clean() {
    rm ${TMP}
}

function fail() {
    clean
    figlet -f big FAIL
    exit 1
}

for PKG in `go list github.com/jollheef/henhouse/... | tr '\n' ' '`; do
    figlet -f big `echo ${PKG} | sed 's|[.a-Z]*/[.a-Z]*/||'`
    echo '---------------' LINT '---------------'
    golint ${PKG} | tee ${TMP} \
        && cat ${TMP} | wc -l | grep 0  >/dev/null && echo "ok" || fail
    echo
    echo '---------------' TEST '---------------'
    go test -v -covermode=count -coverprofile=coverage.out ${PKG} || fail
    echo
    echo
done

clean
figlet -f big OK
