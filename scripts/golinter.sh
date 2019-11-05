#!/bin/bash -e

# Copyright Greg Haskins All Rights Reserved.
#
# SPDX-License-Identifier: Apache-2.0

declare -a vendoredModules=(
"./wasmcc/wasmcc*.go"
"./tools/file-encoder/*.go"
"./integration/e2e/*.go"
)

for i in "${vendoredModules[@]}"
do
    echo ">>>Checking $i with goimports"
    OUTPUT="$(goimports -l ./$i || true )"
    if [[ $OUTPUT ]]; then
        echo "The following files contain goimports errors"
        echo $OUTPUT
        echo "The goimports command 'goimports -l -w' must be run for these files"
        exit 1
    fi
done