#!/bin/sh -x

while getopts "e:a:s:" opt; do
    case "$opt" in
        e) env=${OPTARG};;
        a) aws_account_id=${OPTARG};;
    esac
done

aws ecr get-login-password --region us-west-2 | docker login --username AWS --password-stdin ${aws_account_id}.dkr.ecr.us-west-2.amazonaws.com

docker tag api ${aws_account_id}.dkr.ecr.us-west-2.amazonaws.com/${env}-api:latest
docker tag api ${aws_account_id}.dkr.ecr.us-west-2.amazonaws.com/${env}-api:$(git rev-parse --short HEAD)
docker push ${aws_account_id}.dkr.ecr.us-west-2.amazonaws.com/${env}-api:$(git rev-parse --short HEAD)
docker push ${aws_account_id}.dkr.ecr.us-west-2.amazonaws.com/${env}-api:latest
