#!/bin/bash

if [ "$1" == "--delete" ]; then
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
export GOPATH=$(pwd)/..:$(pwd)/../cisetup:$(pwd)/../go-bindata
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
for f in ../sonarqube/src/main/resources/bayzr/*.xml; do
    fn=$(basename "$f")
    cp $f /etc/bzr.d/$fn
    chmod 644 /etc/bzr.d/$fn
done
cp ../sonarqube/src/main/resources/bayzr/*.xml /etc/bzr.d/

/usr/bin/go get -u github.com/jteeuwen/go-bindata/...
/usr/bin/go get github.com/kisielk/errcheck
/usr/bin/go get github.com/vaughan0/go-ini
/usr/bin/go get github.com/gin-gonic/gin
/usr/bin/go get github.com/gin-gonic/contrib
/usr/bin/go get github.com/elazarl/go-bindata-assetfs

/usr/bin/go get github.com/gorilla/context
/usr/bin/go get github.com/gorilla/securecookie
/usr/bin/go get github.com/garyburd/redigo/redis
/usr/bin/go get github.com/gin-gonic/contrib/sessions
/usr/bin/go get github.com/gin-gonic/contrib/static
/usr/bin/go get github.com/gin-gonic/contrib/renders/multitemplate
/usr/bin/go get github.com/robfig/cron


/usr/bin/go build -o ../bin/go-bindata ../src/github.com/jteeuwen/go-bindata/go-bindata

cd ..
bin/go-bindata -o cisetup/src/data/data.go -pkg data cisetup/src/data cisetup/src/data/css cisetup/src/data/js cisetup/src/data/fonts cisetup/src/data/js/i18n
cd -

/usr/bin/go build -o ../bin/citool ../cisetup/src/main/main.go
