#!/bin/sh

# Run this script: buildah unshare ./build.sh

echo "Build server container using buildah"

BUILDAH_HISTORY=true
NAME=docker.io/mikegolovanov/grpc-users-server
VERSION=1.0 

container=$(buildah from alpine)
echo "Create container $container"
mnt=$(buildah mount $container)
echo "Mount container filesystem to $mnt"
buildah copy $container ./users_server /
buildah umount $container

buildah config --entrypoint "/users_server" $container
buildah config --port 7777 $container 
img=$(buildah commit $container $NAME:$VERSION)

buildah tag $NAME:$VERSION $NAME:lastest
echo "Create image $img done"
 
buildah login -u mikegolovanov docker.io
buildah push $img $NAME:$VERSION 
buildah logout docker.io