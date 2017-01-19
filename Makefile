# This is how we want to name the binary output
BINARY=access_logs_stats

PACKAGE=github.com/blackbass1988/access_logs_stats


# These are the values we want to pass for Version and BuildTime
VERSION=0.5.3
BUILD_TIME=`date +%FT%T%z`

# Setup the -ldflags option for go build here, interpolate the variable values
LDFLAGS=-s -w -X ${PACKAGE}/core.VERSION=${VERSION} -X ${PACKAGE}/core.BUILD_TIME=${BUILD_TIME}

DEBUG= -X github.com/blackbass1988/access_logs_stats/core/debug=1

all: fmt test build

debug: LDFLAGS += ${DEBUG}
debug: fmt test build


fmt:
	go fmt ./...

test:
	go test ./...

build:
	go build -ldflags "${LDFLAGS}" -o ${BINARY} ${PACKAGE}

.PHONY: clean
clean:
	if [ -f ${BINARY} ] ; then rm ${BINARY} ; fi
