TARGET=./bin/panic-middleware
BUILD_CMD="go build -o ${TARGET} main.go"

build:
	eval ${BUILD_CMD}

run-local:
	air -build.cmd ${BUILD_CMD} -build.bin ${TARGET}