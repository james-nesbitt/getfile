.PHONY: all build local vendor fmt test binary install clean

# This script is the switchboard for the make targets
MAKE_SCRIPT="./make/make.sh"

# BY default do the all build
default: all

# Typical make all target
all: build

# Typical build routine
build: vendor test binary install

# Quick developer build
local: fmt binary install
# Full developer build, best to do before Pull Requests
local-full: fmt test vendor local binary


# Run the formatting on all source
fmt:
	${MAKE_SCRIPT} fmt

# Run all go tests
test:
	${MAKE_SCRIPT} test

# Build a binary executable
binary:
	${MAKE_SCRIPT} binary

# Make sure that all the vendor dependencies are available and properly versioned
vendor:
	${MAKE_SCRIPT} vendor

# Remove the compiled binary
clean:
	${MAKE_SCRIPT} clean
