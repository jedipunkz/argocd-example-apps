#/bin/bash

IMAGE_TAG=`git rev-parse --short HEAD`

for dirname in ./bots/*/ ; do
    ECR_REPOSITORY=`echo $dirname | sed 's/\/bots//g' | sed 's/\.//g' | sed 's/\///g'`
    cd $dirname
    docker build -t $ECR_REPOSITORY:$IMAGE_TAG .
    cd ../..
done
