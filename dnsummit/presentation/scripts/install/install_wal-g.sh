#! /bin/bash
rm -rf install
wget https://github.com/wal-g/wal-g/releases/download/v0.2.14/wal-g.linux-amd64.tar.gz
mkdir install && mv wal-g.linux-amd64.tar.gz install
cd install && tar zxvf wal-g.linux-amd64.tar.gz && mv -f wal-g $HOME/bin
cd .. && cp configs/wal-g_minio_conf.json "$HOME/.walg.json"