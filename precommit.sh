#!/bin/sh

TMP=`mktemp`
COVERAGE=coverage.out

MIN_COV=80

function clean() {
    rm ${TMP}
}

function fail() {
    clean
    figlet -f big FAIL
    exit 1
}

rm ${COVERAGE}

for PKG in `go list github.com/jollheef/henhouse/... | tr '\n' ' '`; do
    figlet -f big `echo ${PKG} | sed 's|[.a-Z]*/[.a-Z]*/||'`
    echo '---------------' LINT '---------------'
    golint ${PKG} | tee ${TMP} \
        && cat ${TMP} | wc -l | grep 0  >/dev/null && echo "ok" || fail
    echo
    echo '---------------' TEST '---------------'
    go test -v -covermode=count -coverprofile=${COVERAGE} ${PKG} \
        | tee ${TMP} || fail

    grep 'no test files' ${TMP} >/dev/null || {
        CUR_COV=$(cat ${TMP}  |grep 'coverage: ' | awk '{print $2}' | sed 's/%//')
        python -c "exit(1${CUR_COV} < 1${MIN_COV})" \
            || { echo "Test coverage less than ${MIN_COV}%"; fail; }
    }
    echo
    echo
    cat ${COVERAGE} | grep ' 0$'
    echo
    echo
done

clean
figlet -f big OK
