#!/usr/bin/env bash
####################
#
# Generate a test binary file for testing downloads
#
####

FILENAME="${1:-testfile.bin}""

dd bs=1024 if=/dev/random of="./${FILENAME}" status=progress count=10000
