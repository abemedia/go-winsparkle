#!/bin/sh -e

VERSION="0.8.3"
NAME="WinSparkle-$VERSION"

wget https://github.com/vslavik/winsparkle/releases/download/v$VERSION/$NAME.zip
unzip $NAME.zip
rm -f $NAME.zip

find dll -type f -name "*.dll" -delete
cp -Lr $NAME/Release/WinSparkle.dll dll/x86
cp -Lr $NAME/x64/Release/WinSparkle.dll dll/x64
cp -Lr $NAME/ARM64/Release/WinSparkle.dll dll/arm64
rm -rf $NAME

sed "s/version = \".*\"/version = \"$VERSION\"/g" dll/dll.go > dll/dll.go.tmp
mv dll/dll.go.tmp dll/dll.go
