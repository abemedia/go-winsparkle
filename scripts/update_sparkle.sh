#!/bin/sh

VERSION="0.8.1"
NAME="WinSparkle-$VERSION"

wget https://github.com/vslavik/winsparkle/releases/download/v$VERSION/$NAME.zip
unzip $NAME.zip
rm -f $NAME.zip

find dll -type f -delete
cp -Lr $NAME/Release/WinSparkle.dll dll/x86
cp -Lr $NAME/x64/Release/WinSparkle.dll dll/x64
cp -Lr $NAME/ARM64/Release/WinSparkle.dll dll/arm64
rm -rf $NAME
