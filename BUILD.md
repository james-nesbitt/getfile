# BUILDING Instructions

there are 3 options for building:

## go build it yourself

Make sure that you:

- have go => 1.9
- handle the vendor requirements (git submodules --recursive)
- specify a binary name yourself, otherwise you will get a binary named src


## use the make build target

- you will need go => 1.9 (probably not, but that is what was tested)
- use `make build`

- the binary will be ./bin/getfile

## Use the build.sh script for a dockerized compile.

- you will need docker running
- it will pull a docker image (docker permissions required)
- it will run a docker container (docker permissions required)

- the binary will be ./bin/getfile

## Other options considered

I considered a Dockerfile build, but that seemed superfluous.  It would make sense for platforms that are not
available for go binaries.
