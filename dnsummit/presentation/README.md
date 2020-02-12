# DNet Summit 2020 Presentation

# Slide Setup

1) Install cleaver with `npm i -g cleaver`
2) Build the latest version of the sldies with `make build-workshop`
3) Open the slides with the brave browser using `make open-workshop` or manually open the `slides/WORKSHOP-CLEAVER_5-cleaver.html` file. 

# Demo

## Install

To setup the demo you will need a valid installation of docker-compose, and a locally available TemporalX docker image. Unfortunately these are hard to come by so you'll need to use our remote development node that sits at https://xapi-dev.temporal.cloud but its behind an extremely throttled network connection.

1) Install `rclone` with `make install-rclone`
2) Install `wal-g` with `make install wal-g`
3) Install `temporal` with `make install-temporal
4) Install `mc` the minio command line client
5) Setup `mc` with the following configuration
```
s3x-local
  URL       : http://127.0.0.1:9000
  AccessKey : minio
  SecretKey : minio123
  API       : s3v2
  Lookup    : aut
```
6) Start the test environment with `make testenv`


## Rclone Demo

The rclone demo show cases using rclone to backup your filesystem to s3, subsequently storing your filesystem on ipfs. 

To view the rclone demo, run `make rclone-demo`.

## Database Demo

The databse demo show cases backing up your postgresql database to s3, subsequently storing your database backups on ipfs.

To view the database demo, run `make database-demo`
