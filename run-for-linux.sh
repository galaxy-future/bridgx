#!/usr/bin/env bash
# deploy mysql
docker run -d --name bridgx_db -e MYSQL_ROOT_PASSWORD=mtQ8chN2 -e MYSQL_DATABASE=bridgx -e MYSQL_USER=gf -e MYSQL_PASSWORD=db@galaxy-future.com -p 3306:3306 -v $(pwd)/init/mysql:/docker-entrypoint-initdb.d yobasystems/alpine-mariadb:10.5.11
# deploy etcd
docker run -d --name bridgx_etcd -e ALLOW_NONE_AUTHENTICATION=yes -e ETCD_ADVERTISE_CLIENT_URLS=http://etcd:2379 -p 2379:2379 -p 2380:2380 bitnami/etcd:3
# deploy api
docker run -d --name bridgx_api --network host -v $(pwd)/conf/config.yml.prod:/home/tiger/api/conf/config.yml.prod galaxyfuture/bridgx-api:v0.1.0 bin/wait-for-api.sh
# deploy sheduler
docker run -d --name bridgx_scheduler --network host -v $(pwd)/conf/config.yml.prod:/home/tiger/app/conf/config.yml.prod  galaxyfuture/bridgx-scheduler:v0.1.0 bin/wait-for-scheduler.sh
