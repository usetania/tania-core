#!/bin/bash

set -e

echo "Downloading osx cross..."
wget -q https://s3-ap-southeast-1.amazonaws.com/tanibox-build-tools/osx-cross.10.11.sdk.txz

echo "Extracting..."

CURRENT_DIR=$(pwd)
mkdir ~/osx_tools && \
  cd ~/osx_tools && \
  tar xJf ${CURRENT_DIR}/osx-cross.10.11.sdk.txz && \
  cd ${CURRENT_DIR}


