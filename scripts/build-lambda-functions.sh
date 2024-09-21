#!/bin/bash

set +e

for d in cmd/aws/lambda/* ; do
    GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -mod=readonly -o "$d/bootstrap" "./$d"

    zip "lambda-function-$(basename $d).zip" "$d/bootstrap"

    rm "$d/bootstrap"
done
