#!/usr/bin/env bash

read -r -p "Are you sure to install mysql? [Y/N] " input

case $input in
    [yY][eE][sS]|[yY])
		echo "Installing mysql..."
    # deploy mysql
    docker run -d --name bridgx_db -e MYSQL_ROOT_PASSWORD=mtQ8chN2 -e MYSQL_DATABASE=bridgx -e MYSQL_USER=gf -e MYSQL_PASSWORD=db@galaxy-future.com -p 3306:3306 -v $(pwd)/init/mysql:/docker-entrypoint-initdb.d yobasystems/alpine-mariadb:10.5.11
		;;

    [nN][oO]|[nN])
		echo "Skip mysql install, please check conf/config.yml mysql config, and import init/mysql/* to existing mysql for first install."
    ;;
    *)
		echo "Invalid input..."
		exit 1
		;;
esac

read -r -p "Are you sure to install etcd? [Y/N] " input

case $input in
    [yY][eE][sS]|[yY])
		echo "Installing etcd..."
    # deploy etcd
    docker run -d --name bridgx_etcd -e ALLOW_NONE_AUTHENTICATION=yes -e ETCD_ADVERTISE_CLIENT_URLS=http://etcd:2379 -p 2379:2379 -p 2380:2380 bitnami/etcd:3
		;;

    [nN][oO]|[nN])
		echo "Skip etcd install, please check conf/config.yml etcd config."
    ;;
    *)
		echo "Invalid input..."
		exit 1
		;;
esac

# deploy api
docker run -d --name bridgx_api --network host -v $(pwd)/conf/config.yml.prod:/home/tiger/api/conf/config.yml.prod galaxyfuture/bridgx-api:"${1:-latest}" bin/wait-for-api.sh
# deploy scheduler
docker run -d --name bridgx_scheduler --network host -v $(pwd)/conf/config.yml.prod:/home/tiger/scheduler/conf/config.yml.prod  galaxyfuture/bridgx-scheduler:"${1:-latest}" bin/wait-for-scheduler.sh
