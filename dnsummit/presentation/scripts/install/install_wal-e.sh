#! /bin/bash

cd wal-e
sudo python3 setup.py install
if [[ "$?" != "0" ]]; then
    echo "[ERROR] failed to install wal-e"
    exit 1
fi
wal-e --help
if [[ "$?" !=  "0" ]]; then
    echo "[ERROR] failed to run wal-e --help"
    exit 1
fi
echo "[INFO] installed wal-e"