#!/bib/bash

echo "Welcome to bayzr interactive build preparator"
echo "Enter build version and release number"
read -p "Version: " -e version
read -p "Release: " -e release
VER="$version-$release"
if [ -z "$VER" ]; then
	echo "Version and release can't be empty"
	exit 255
fi

PRP_DIR="bayzr-$version"
rm -rf "$PRP_DIR"
mkdir -p "$PRP_DIR/bin"
cp -R ../cfg "$PRP_DIR"
cp -R ../src "$PRP_DIR"
mkdir -p "$PRP_DIR/rpm"

cp ../rpm/COPYRIGHT "$PRP_DIR/rpm/"
cp ../rpm/gpl-3.0.txt "$PRP_DIR/rpm/"

mv "$PRP_DIR/src/main/main.go" "$PRP_DIR/src/main/main.go.bak"

sed "s/const APP_VERSION = \"0.1\"/const APP_VERSION = \"$VER\"/g" "$PRP_DIR/src/main/main.go.bak" > "$PRP_DIR/src/main/main.go"

rm -rf "$PRP_DIR/src/main/main.go.bak"

tar zcvf "$PRP_DIR.tar.gz" "$PRP_DIR"

sed "s/VERRR/$version/g" bzr-analyzer.spec.tpl | sed "s/RELLL/$release/g" > bzr-analyzer.spec
cat changelog.list >> bzr-analyzer.spec

rm -rf "$PRP_DIR"




