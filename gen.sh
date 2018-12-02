#!/bin/bash

HASH=`git rev-parse HEAD`

cat << EOF > version_linux.go
package main

//go:generate bash gen.sh
const commitHash = "$HASH"
EOF


cat << EOF > version_windows.go
package main

//go:generate powershell .\gen.ps1
const commitHash = "${HASH}"
EOF