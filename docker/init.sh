#!/bin/bash
set -x


function do_root {

   apt-get update
   apt-get install -y --force-yes gcc runit wget git

   adduser --disabled-login --gecos 'app' app

   cd /srv/install
   cp godeb-helper /usr/bin/godeb-helper

   godeb-helper
   cp start-app /usr/bin/start-app

   # Run as regular user now
   cp -v $0 /tmp/init.sh
   chown app /tmp/init.sh
   su -s /bin/bash - app -c "bash /tmp/init.sh"
   exit $?
}

function do_non_root {
   echo $GOPATH
   cd
   pwd
   source .bashrc
   go get -u github.com/sigmonsays/paste/pasted
   exit 0
}

WHOAMI="$(whoami)"
if [ "$WHOAMI" == "root" ] ; then
   do_root $@
else
   do_non_root $@
fi

