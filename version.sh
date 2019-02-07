#!/bin/sh -e

VFILE=version.go

echo "// Updated automatically (altered manually just prior to each release)" > $VFILE
echo "" >> $VFILE
echo "package sqlapi" >> $VFILE
echo "" >> $VFILE
echo "const Version = \"$(git describe --tags --always 2>/dev/null)\"" >> $VFILE
