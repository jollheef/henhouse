#!/bin/sh

for PKG in `go list github.com/jollheef/henhouse/... | tr '\n' ' '`; do
    figlet -f big `echo ${PKG} | sed 's|[.a-Z]*/[.a-Z]*/||'`
    echo '---------------' LINT '---------------'
    golint ${PKG}
    echo
    echo '---------------' TEST '---------------'
    go test -v -covermode=count -coverprofile=coverage.out ${PKG} || exit 1
    echo
    echo
done
