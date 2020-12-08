#/bin/bash

IMAGE_TAG=`git rev-parse --short HEAD`
ECR_REGISTRY=$1
ENV=$2

for dirname in ./bots/*/ ; do
    ECR_REPOSITORY=`echo $dirname | sed 's/\/bots//g' | sed 's/\.//g' | sed 's/\///g'`
    cd $dirname
    docker tag $ECR_REPOSITORY:$IMAGE_TAG $ECR_REGISTRY/infra-bot-$ENV/$ECR_REPOSITORY:$IMAGE_TAG
    docker push $ECR_REGISTRY/infra-bot-$ENV/$ECR_REPOSITORY:$IMAGE_TAG
    docker tag $ECR_REPOSITORY:$IMAGE_TAG $ECR_REGISTRY/infra-bot-$ENV/$ECR_REPOSITORY:latest
    docker push $ECR_REGISTRY/infra-bot-$ENV/$ECR_REPOSITORY:latest

    kustomize edit set image $ECR_REPOSITORY=$ECR_REGSTRY/infra-bot-$ENV/$ECR_REPOSITORY:$IMAGE_TAG
    cd ../..
done
