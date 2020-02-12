#! /bin/bash

rm -rf install && mkdir install
cd install && wget -O minio https://gateway.temporal.cloud/ipfs/QmbfHK3QaXU6vb5Rdvzb4CdPVGhAaLBgsYveYtB7VWyPJU
cp minio $HOME/bin && chmod a+x $HOME/bin/minio