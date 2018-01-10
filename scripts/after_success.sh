#!/bin/bash

set -e

TIMESTAMP="$(date +%Y%m%d-%s)-b${TRAVIS_BUILD_NUMBER}"

echo -n "Archiving..."

tar cvJf tania-server-${TIMESTAMP}-linux-arm.txz terra.linux.arm public
tar cvJf tania-server-${TIMESTAMP}-linux-amd64.txz terra.linux.amd64 public
7za a -t7z -mx=9 tania-server-${TIMESTAMP}-windows-amd64.7z terra.windows.amd64.exe public


if  [ "x$TRAVIS_PULL_REQUEST" == "xfalse" ] 
then

  echo -n "Uploading to S3..."

  aws s3 cp tania-server-${TIMESTAMP}-linux-arm.txz s3://tanibox-terra/archives/ --storage-class STANDARD_IA
  aws s3 cp tania-server-${TIMESTAMP}-linux-amd64.txz s3://tanibox-terra/archives/ --storage-class STANDARD_IA
  aws s3 cp tania-server-${TIMESTAMP}-windows-amd64.7z s3://tanibox-terra/archives/ --storage-class STANDARD_IA

fi
