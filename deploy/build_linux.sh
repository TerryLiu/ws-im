#!/usr/bin/env bash

appName=wsim

if [ $# -eq 0 ];
then nextVersion=$(./increment_version.sh -p ${BASH_REMATCH[2]});
else nextVersion=$(./increment_version.sh $1 ${BASH_REMATCH[2]});
fi

echo "next version"
echo $nextVersion;

rm -rf $appName
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o $appName *.go

docker build -t hub.docker.com/cqliuz/$appName:$nextVersion .
docker push hub.docker.com/cqliuz/$appName:$nextVersion

git add ../
git commit -m "$nextVersion $1"
git push
date