#!/bin/sh
START_TIME=`date +%s`

RUNDIR=$(pwd)

cd $(dirname ${0})
COMMIT_ID=`git --no-pager log --format="%H" -n 1`

cd ${RUNDIR}

BUILD_DATE=`date -u +%d.%m.%Y`
BUILD_TIME=`date -u +%H:%M:%S`

LDFLAGS="-X main.CommitID=${COMMIT_ID}"
LDFLAGS+=" -X main.BuildDate=${BUILD_DATE}"
LDFLAGS+=" -X main.BuildTime=${BUILD_TIME}"

mkdir -p ${GOPATH}/bin

go build -ldflags "${LDFLAGS}" -o ${GOPATH}/bin/henhouse github.com/jollheef/henhouse
go build -ldflags "${LDFLAGS}" -o ${GOPATH}/bin/henhousectl github.com/jollheef/henhouse/cli/henhousectl

END_TIME=`date +%s`
RUN_TIME=$((END_TIME-START_TIME))
echo "Build done in ${RUN_TIME}s"
