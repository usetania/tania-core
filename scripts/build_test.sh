#/bin/bash

set -e

# Running any go tests
echo "Running go test.."
go test ./...

echo "Building golang binaries..."
# Build and test Golang
make linux-arm linux-amd64 windows osxcross

echo "Configuring binary for running ..."
# Setting up configuration
cp conf.json.example conf.json
sed -i.bak "s|/Users/user/Code/golang|$TRAVIS_BUILD_DIR/gopath|g" conf.json

echo "Starting server for E2E testing ..."

# Run golang on linux
./terra.linux.amd64 > /dev/null 2>&1 &
TERRA_PID=$!

echo "Server has running running in the background at pid ${TERRA_PID}"

echo "Running Front-End Unit tests ..."
# build and run unit test
yarn && yarn run unit

echo "Running end to end tests ..."
# build and test e2e
yarn run production && yarn run cypress:run

echo "Killing Server [$TERRA_PID] ..."

# Move the screenshoot and recorded video from the test result into public folder
mkdir public/assets
cp -rf resources/tests/assets/videos public/assets/
cp -rf resources/tests/assets/screenshots public/assets/

kill -s TERM $TERRA_PID
