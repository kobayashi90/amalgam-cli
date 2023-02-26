#!/bin/bash

cmdVersion=$(cat cmd/cmd.go | grep app.Version | cut -d'=' -f2 | tr -d ' "')
makeVersion=$(cat makefile | grep "^VERSION" | cut -d'=' -f2 | tr -d ' "')

if [ "$cmdVersion" != "$makeVersion" ]; then
  echo -e "Version Mismatch:\n\tcmd: $cmdVersion\n\tmakefile: $makeVersion"
  exit 1
fi