#! /bin/bash

PGHOST=127.0.0.1:5432
PGUSER=postgres
PGPASSWORD=password123
DBDIR=$(docker inspect dnet_summit_2020_workshop_postgres_1 | grep "Source" | tr -d '",' | awk '{print $2}')
export PGHOST
export PGUSER
export PGPASSWORD

echo "[INFO] running database migrations to generate test data"
make migrate-db
echo "[INFO] copying wal-g config"
sudo cp ./configs/wal-g_minio_conf.json /root/.walg.json
echo "[INFO] backing up postgres database"
sudo -E wal-g backup-push "$DBDIR"
echo "database directory: $DBDIR"
