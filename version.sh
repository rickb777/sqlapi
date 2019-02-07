#!/bin/sh -e

VFILE=util/version.go

echo "// Updated automatically" > $VFILE
echo "" >> $VFILE
echo "package util" >> $VFILE
echo "" >> $VFILE
echo "const Version = \"$(git describe --tags 2>/dev/null)\"" >> $VFILE
