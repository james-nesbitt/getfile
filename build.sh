#!/usr/bin/env bash

#
# Build binary in a container
#
# @NOTE to specify a different os/arch:
#    - GOOS : linux darwin windows
#    - GOARCH : amd64 arm arm64
#

source make/.os-detect
source make/config.sh

EXEC_PATH="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
INTERNAL_LIBRARY_PATH="github.com/james-nesbitt/getfile"

echo "***** Building .

This will build as a binary for '$GOOS-$GOARCH'.

(Override this by setting \$GOOS and \$GOARCH environment variables)

 **** Building in containerized golang environment
 "

# some sanity stuff, to prevent docker related permissions issues
mkdir -p ".git/modules/vendor"
chmod u+x Makefile

# Run the build inside a container
#
#  - volumify the submodule changes
#  - build in a valid gopath to get active vendor dependencies
#  - pass in env variables for environment control
docker run --rm -ti \
	-v "${EXEC_PATH}:/go/src/${INTERNAL_LIBRARY_PATH}" \
	-v "/go/src/${INTERNAL_LIBRARY_PATH}/.git/modules/vendor" \
	-v "/go/src/${INTERNAL_LIBRARY_PATH}/vendor" \
	-e "GOOS=${GOOS}" \
	-w "/go/src/${INTERNAL_LIBRARY_PATH}" \
	golang:${GOVERSION} \
	make build

echo "

Finished building the application inside the container.  If an error occured
during the golang compile, then you would have seen it reported above.

"

echo " **** Containerized build complete

an executable binary has (hopefully) now been built
in ${BUILD_BINARY_PATH}

"
