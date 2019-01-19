#/bin/bash

set -e

# Running any go tests
echo "Running go test.."
go test ./...

# Create empty sqlite db file. Will be overwritten everytime.
echo "Create empty sqlite db file"
touch db/sqlite/tania.db
chmod 775 db/sqlite/tania.db

echo "Building golang binaries..."
# Build and test Golang
make linux-arm linux-amd64 windows osxcross

echo "Configuring binary for running ..."
# Setting up configuration
cp conf.json.example conf.json
sed -i.bak "s|/Users/user/Code/golang/src/github.com/Tanibox/tania-server|$TRAVIS_BUILD_DIR|g" conf.json
# Set DEMO_MODE to true to turn off the token validation
DEMO_MODE=false
echo "DEMO_MODE is set to ${DEMO_MODE}"

echo "Starting server for E2E testing ..."

# Run golang on linux
./tania.linux.amd64 > /dev/null 2>&1 &
TANIA_PID=$!

echo "Server has running running in the background at pid ${TANIA_PID}"

echo "Running Front-End Unit tests ..."
# build and run unit test
yarn && yarn run unit

echo "Running end to end tests ..."
# build and test e2e
yarn run production && yarn run cypress:run

echo "Killing Server [$TANIA_PID] ..."

# Move the screenshoot and recorded video from the test result into public folder
mkdir public/assets
cp -rf resources/tests/assets/videos public/assets/

kill -s TERM $TANIA_PID
