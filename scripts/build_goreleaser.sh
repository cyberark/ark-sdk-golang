#!/bin/bash

if [ -z "$BUILD_VERSION" ]
then
	export BUILD_VERSION="v1.0.0"
fi

if [ -z "$BUILD_NUMBER" ]
then
	export BUILD_NUMBER="0"
fi

SCRIPTPATH=$(dirname "$0")
OUT_PATH=$SCRIPTPATH/../dist

function check_if_exists() {
  fileName="$1"
  if [ ! -f "$OUT_PATH/$fileName" ]
  then
     echo "$fileName does not exists"
     exit 1
  fi
}

CURRENT_BRANCH=$(git rev-parse --abbrev-ref HEAD)
if [[ "$CURRENT_BRANCH" != "main" && "$CURRENT_BRANCH" != "master" ]]
then
  goreleaser build --clean --snapshot
else
  goreleaser build --clean
fi

check_if_exists "unix_darwin_amd64_v1/ark-darwin"
check_if_exists "unix_darwin_arm64_v8.0/ark-darwin"
check_if_exists "unix_freebsd_386_sse2/ark-freebsd"
check_if_exists "unix_freebsd_amd64_v1/ark-freebsd"
check_if_exists "unix_freebsd_arm64_v8.0/ark-freebsd"
check_if_exists "unix_freebsd_arm_6/ark-freebsd"
check_if_exists "unix_linux_386_sse2/ark-linux"
check_if_exists "unix_linux_amd64_v1/ark-linux"
check_if_exists "unix_linux_arm64_v8.0/ark-linux"
check_if_exists "unix_linux_arm_6/ark-linux"
check_if_exists "win_windows_386_sse2/ark-windows.exe"
check_if_exists "win_windows_amd64_v1/ark-windows.exe"
check_if_exists "win_windows_arm64_v8.0/ark-windows.exe"
