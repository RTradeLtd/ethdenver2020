.PHONY: open-workshop
open-workshop: build-workshop
	(cd slides ; brave-browser WORKSHOP-CLEAVER_5-short-cleaver.html)
	
.PHONY: build-workshop
build-workshop:
	(cd slides ; cleaver WORKSHOP-CLEAVER_5-short.md)

.PHONY: install-rclone
insatll-rclone:
	./scripts/install/install_rclone.sh

.PHONY: install-wal-g
install-wal-g:
	./scripts/install/install_wal-g.sh

.PHONY: install-temporal
install-temporal:
	./scripts/install/install_temporal.sh

.PHONY: testenv
testenv:
	docker-compose -f docker-compose.yml up -d
	mc mb s3x-local/testdata

.PHONY: migrate-db
migrate-db:
	temporal -config configs/temporal_database.json --db.no_ssl migrate

.PHONY: rclone-demo
rclone-demo:
	./scripts/demo/rclone_backup_demo.sh

.PHONY: database-demo
database-demo:
	./scripts/demo/database_backup_demo.sh