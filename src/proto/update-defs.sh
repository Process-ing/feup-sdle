#!/bin/sh

ROOT_DIR=$(dirname "$0")/..

# Generate JavaScript/TypeScript code from the proto files
pbjs -t static-module -w es6 -o "$ROOT_DIR/client/lib/proto/global.js" "$ROOT_DIR/proto/global.proto"
pbts -o "$ROOT_DIR/client/lib/proto/global.d.ts" "$ROOT_DIR/client/lib/proto/global.js"