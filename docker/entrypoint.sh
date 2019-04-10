#!/bin/sh
export AWS_DEFAULT_REGION=eu-central-1
aws s3 cp s3://${BUCKET_NAME}/config.yml /app/config.yml
aws s3 cp s3://${BUCKET_NAME}/supervisord.conf /etc/supervisor/conf.d/supervisord.conf

/usr/bin/supervisord