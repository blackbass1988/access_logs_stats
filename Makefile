# This is how we want to name the binary output
BINARY=access_logs_stats

# These are the values we want to pass for Version and BuildTime
BUILD_TIME=`date +%FT%T%z`
COMMIT=`git rev-parse HEAD`
BRANCH=`git rev-parse --abbrev-ref HEAD`

# Setup the -ldflags option for go build here, interpolate the variable values
#LDFLAGS=-s -w -X main.buildTime=${BUILD_TIME}
LDFLAGS=-X main.buildTime=${BUILD_TIME} -X main.commit=${COMMIT} -X main.branch=${BRANCH}

DEBUG= -X github.com/blackbass1988/access_logs_stats/core/debug=1

all: fmt test build

debug: LDFLAGS += ${DEBUG}
debug: fmt test build


fmt:
	go fmt ./...

test:
	go test ./core/...

build:
	go build -ldflags "${LDFLAGS}" ./cmd/...

.PHONY: clean
clean:
	if [ -f ${BINARY} ] ; then rm ${BINARY} ; fi
