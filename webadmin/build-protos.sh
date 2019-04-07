#!/bin/bash

OUTPUT=./src/radioserver

protoc \
    -I../ \
    --plugin=protoc-gen-ts=./node_modules/.bin/protoc-gen-ts \
    --js_out=import_style=commonjs,binary:$OUTPUT \
    --ts_out=service=true:$OUTPUT \
    ../protocol/*.proto

# Fix bug with create-react-app fucking eslint ignoring .eslintignore

cd src/radioserver/protocol
for i in *_pb*.js
do
    echo "Fixing $i"
    sed -i '1i/* eslint-disable */' $i
    sed -i 's/var jspb/\/\/ @ts-ignore\nvar jspb/g' $i
done
