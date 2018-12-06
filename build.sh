#!/bin/sh

cd $(dirname $BASH_SOURCE)
DIR=$(pwd)

if [ -z "$(git status --porcelain)" ]; then
  # Working directory clean
  echo
else
  echo "Project must be clean before you can build"
  exit 1
fi

rm -rf build

mkdir -p build/darwin
mkdir -p build/linux

OLD_VER=$(awk '/^const version =/ {print $0}' version.go)
NEW_VER=$(echo $OLD_VER | awk -F'[ .]' '/^const version =/ {print $1,$2,$3,$4"."$5"."$6+1"\""}')

sed -i .bak "s/$OLD_VER/$NEW_VER/g" version.go

VER=$(awk -F'[ ."]' '/^const version =/ {print $5"."$6"."$7}' version.go)

git add version.go

git commit -m "New version - $VER"

GIT_COMMIT=$(git rev-parse HEAD)

env CGO_ENABLED=1 go build -o build/darwin/blocksaver_$VER -ldflags="-X main.commit=${GIT_COMMIT}"

# Because ZMQ uses C bindings, we cannot cross compile.  Instead use Docker...
tar xvfz zeromq-4.1.4.tar.gz
docker build -t maestro .
docker run -v $DIR:/paymaster -i maestro <<EOL
cd /paymaster/zeromq-4.1.4
apt-get install -y libsodium-dev
./configure
make
make install
cd /paymaster
go get -v ./...
CGO_ENABLED=1 go build -a -o build/linux/blocksaver_$VER  -ldflags="-X main.commit=${GIT_COMMIT}"
exit
EOL

rm -rf zeromq-4.1.4/

cp settings.conf start.sh stop.sh tailLog.sh build/darwin
cp settings.conf start.sh stop.sh tailLog.sh paymaster.service build/linux
