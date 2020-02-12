#! /bin/bash


AWS_ACCESS_KEY_ID="minio"
export AWS_ACCESS_KEY_ID
AWS_SECRET_ACCESS_KEY="minio123"
export AWS_SECRET_ACCESS_KEY
WALG_S3_PREFIX="s3://testdata"
export WALG_s3_PREFIX
AWS_ENDPOINT="http://206.116.153.42:25001"
export AWS_ENDPOINT
AWS_S3_FORCE_PATH_STYLE="true" # we dont support subdomain style addressing
export AWS_S3_FORCE_PATH_STYLE
AWS_REGION=us-east-1 # not required, but here for brevity
export AWS_REGION
# WALG_S3_CA_CERT_FILE="/path/to/custom/ca/file"
# export WALG_S3_CA_CERT_FILE