#!/bin/bash

if [ -z "$GIT_COMMIT" ]
then
	export GIT_COMMIT=$(git rev-list -1 HEAD)
fi

if [ -z "$BUILD_DATE" ]
then
	export BUILD_DATE=$(date '+%d-%m-%Y %T')
fi

if [ -z "$BUILD_VERSION" ]
then
	export BUILD_VERSION="v1.0.0"
fi

if [ -z "$BUILD_NUMBER" ]
then
	export BUILD_NUMBER="0"
fi

SCRIPTPATH=$(dirname "$0")

IN_PATH=$SCRIPTPATH/../cmd/ark/.
OUT_PATH=$SCRIPTPATH/../bin

function check_if_exists() {
  fileName="$1"
  if [ ! -f "$OUT_PATH/$fileName" ]
  then
     echo "$fileName does not exists"
     exit 1
  fi
}
UNIX_LD_FLAGS="-s -w -X 'main.GitCommit=$GIT_COMMIT' -X 'main.BuildDate=$BUILD_DATE' -X 'main.Version=$BUILD_VERSION' -X 'main.BuildNumber=$BUILD_NUMBER'"
WINDOWS_LD_FLAGS="-X 'main.GitCommit=$GIT_COMMIT' -X 'main.BuildDate=$BUILD_DATE' -X 'main.Version=$BUILD_VERSION' -X 'main.BuildNumber=$BUILD_NUMBER'"
env GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 go build -ldflags "$UNIX_LD_FLAGS" -o "$OUT_PATH/ark-darwin" "$IN_PATH/ark.go"
env GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -ldflags "$WINDOWS_LD_FLAGS" -o "$OUT_PATH/ark-windows.exe" "$IN_PATH/ark.go"
env GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags "$UNIX_LD_FLAGS" -o "$OUT_PATH/ark-linux" "$IN_PATH/ark.go"

echo "********* $OUT_PATH ********"
ls -la "$OUT_PATH"
echo "*****************"

check_if_exists "ark-darwin"
check_if_exists "ark-windows.exe"
check_if_exists "ark-linux"

exit 0
