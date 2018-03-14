#!/usr/bin/env bash
####################
#
# A switchboard used to setup make steps
#
# - configure GO building environment
# - config some ENVs for GO build definition
# - config some ENV variables for the app itself
#
####

BUILD_PATH="./bin"

# Go build config
# @TODO We should be determining these automatically somehow?
export GOOS="${GOOS:-linux}" # Perhaps you would prefer "osx" ?
export GOARCH="${GOARCH:-amd64}"
export GOVERSION="1.9"

# Application locations (including custom vendors that may need test run
GO_PKG="getfile"
GO_TARGET="${GO_TARGET:-src}"
