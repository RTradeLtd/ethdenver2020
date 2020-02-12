#! /bin/bash

DIR=$(apg | xargs -0 | sha256sum | awk '{print $1}')
rclone mkdir s3x-ipfs:"$DIR"
echo "[INFO] running regular copy"
rclone copy ./configs "s3x-ipfs:$DIR" --s3-access-key-id minio --s3-secret-access-key minio123 --s3-upload-cutoff 1M
echo "[INFO] running a 'delta sync'"
rclone sync ./scripts "s3x-ipfs:$DIR" --s3-access-key-id minio --s3-secret-access-key minio123 --s3-upload-cutoff 1M
echo "bucket name: $DIR"