#!/bin/bash

NOW=$(date +%Y%m%d%H%M)
IMG_NAME=operator:${NOW}

ORIGINDIR=$(pwd)

# build image
cp ./release/operator ./dockerfiles/operator/
cd ./dockerfiles/operator/
docker build -t ${IMG_NAME} .
rm ./operator
cd ${ORIGINDIR}

# push image
docker login ${QiniuRegistry} -u ${QiniuRegistryAK} -p ${QiniuRegistrySK} ${QiniuRegistry}
docker tag ${IMG_NAME} ${QiniuRegistryPrefix}/${IMG_NAME}
docker push ${QiniuRegistryPrefix}/${IMG_NAME}