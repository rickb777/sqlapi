#!/bin/sh -e
V=$1
if [ -z "$V" ]; then
  V=$(git describe --tags --always 2>/dev/null)
fi

VFILE=version.go

echo "// Updated automatically (altered manually just prior to each release)" > $VFILE
echo "" >> $VFILE
echo "package sqlapi" >> $VFILE
echo "" >> $VFILE
echo "const Version = \"$V\"" >> $VFILE
