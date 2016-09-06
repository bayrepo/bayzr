#!/bin/bash

if [ "$1"=="--delete" ]; then
    if [ -e /etc/bzr.conf ]; then
        rm -f /etc/bzr.conf
    fi
    if [ -e /etc/bzr.d ]; then
        rm -rf /etc/bzr.d
    fi
    if [ -e /usr/bin/bayzr ]; then
        rm -f /usr/bin/bayzr
    fi
    exit
fi

export GOROOT=/usr/lib/golang/
export PATH=$PATH:$GOROOT/bin
export GOPATH=$(pwd)/..
export PATH=$PATH:$GOPATH/bin

/usr/bin/go get github.com/jroimartin/gocui
/usr/bin/go get github.com/nsf/termbox-go
/usr/bin/go get github.com/mattn/go-runewidth
/usr/bin/go get github.com/go-sql-driver/mysql

/usr/bin/go build -o ../bin/bayzr main

mkdir -p /etc/bzr.d

cp ../bin/bayzr /usr/bin/
chmod 755 /usr/bin/bayzr
cp ../cfg/bzr.conf /etc/
chmod 644 /etc/bzr.conf
for f in ../cfg/*.conf; do
    fn=$(basename "$f")
    if [ "$fn" != "bzr.conf" ]; then
        cp $f /etc/bzr.d/$fn
        chmod 644 /etc/bzr.d/$fn
    fi
done
for f in ../cfg/*.tpl; do
    fn=$(basename "$f")
    if [ "$fn" != "bzr.conf.tpl" -a "$fn" != "checkerplugin.cfd.tpl" ]; then
        cp $f /etc/bzr.d/$fn
        chmod 644 /etc/bzr.d/$fn
    fi
done