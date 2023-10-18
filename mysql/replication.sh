#!/bin/bash

mysqldump -h db -uroot -proot --databases db > dbdump.db
mysql -h replica -uroot -proot < dbdump.db
mysql -h db -uroot -proot < /docker-entrypoint-initdb.d/primary.sql
mysql -h replica -uroot -proot < /docker-entrypoint-initdb.d/replica.sql
