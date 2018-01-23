#/bin/bash

set -e

# Build and test Golang
make linux-arm linux-amd64 windows

CURRENT_DIR=$(pwd)

# Run golang on linux
./terra.linux.amd64 > /dev/null 2>&1 &
TERRA_PID=$!

echo "Terra running in the background at pid ${TERRA_PID}"

# build and run unit test
yarn && yarn run unit

# Setting up configuration
cp conf.json.example conf.json
sed -i.bak "s|/Users/user/Code/golang|$TRAVIS_BUILD_DIR/gopath|g" conf.json

# build and test e2e
yarn run production && yarn run cypress:run 

cd ${CURRENT_DIR}

kill -s TERM $TERRA_PID
