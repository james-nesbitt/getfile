#!/usr/bin/env bash
####################
#
# Make switchboard includer.  Includes a subfile based on a passed target,
# after including a setup/config script.
#
# @NOTE these scripts expect to be run from the project top path, as they
#    use paths relative to the root.
#
# @NOTE Do not do any logic or functionality in this file
#    as it may in some circumstances be sources in an
#    escalated permission environment
#
###

# Include our configure script
source make/.os-detect
source make/config.sh

if [ -z "$GOPATH" ]; then
    echo "WARNING: No GOPATH exists in your environment.  Certain components such as TESTs may produce weird errors"
fi

# Build a bundle
bundle() {
 local bundle="$1"; shift
 echo "---> Make-bundle: $(basename "$bundle") (in $DEST)"
 source "make/$bundle" "$@"
}

if [ $# -gt 0 ]; then
 bundles=($@)
 for bundle in ${bundles[@]}; do
     export DEST=.
     ABS_DEST="$(cd "$DEST" && pwd -P)"
     bundle "$bundle"
     echo
 done
fi
