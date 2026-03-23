#!/usr/bin/env sh
set -eu

ROOT_DIR=$(CDPATH= cd -- "$(dirname "$0")/.." && pwd)
COMPOSE_FILE=docker-compose.infra.dev.yml
APP_CONFIG=${APP_CONFIG:-configs/config.local.yaml}

cd "$ROOT_DIR"

wait_for_health() {
	container_name=$1
	max_attempts=$2
	attempt=1

	while [ "$attempt" -le "$max_attempts" ]; do
		status=$(docker inspect -f '{{if .State.Health}}{{.State.Health.Status}}{{else}}{{.State.Status}}{{end}}' "$container_name" 2>/dev/null || true)
		if [ "$status" = "healthy" ] || [ "$status" = "running" ]; then
			return 0
		fi

		echo "waiting for $container_name to become ready ($attempt/$max_attempts)..."
		sleep 2
		attempt=$((attempt + 1))
	done

	echo "$container_name did not become ready in time" >&2
	return 1
}

read_mysql_config_value() {
	key=$1
	awk -v target="$key" '
		$1 == "mysql:" { in_mysql = 1; next }
		in_mysql && /^[^[:space:]]/ { in_mysql = 0 }
		in_mysql {
			gsub(":", "", $1)
			if ($1 == target) {
				print $2
				exit
			}
		}
	' "$APP_CONFIG"
}

ensure_mysql_database() {
	mysql_user=$(read_mysql_config_value user)
	mysql_password=$(read_mysql_config_value password)
	mysql_database=$(read_mysql_config_value dbname)
	mysql_charset=$(read_mysql_config_value charset)

	if [ -z "$mysql_user" ] || [ -z "$mysql_password" ] || [ -z "$mysql_database" ] || [ -z "$mysql_charset" ]; then
		echo "mysql configuration is incomplete in $APP_CONFIG" >&2
		return 1
	fi

	docker exec sys-admin-serve-mysql-dev \
		mysql -u"$mysql_user" -p"$mysql_password" \
		-e "CREATE DATABASE IF NOT EXISTS \`$mysql_database\` CHARACTER SET $mysql_charset COLLATE utf8mb4_unicode_ci;"
}

docker compose -f "$COMPOSE_FILE" up -d
wait_for_health sys-admin-serve-mysql-dev 30
wait_for_health sys-admin-serve-redis-dev 15
ensure_mysql_database
APP_CONFIG="$APP_CONFIG" go run ./cmd/migrate up
APP_CONFIG="$APP_CONFIG" go run ./cmd/seed
APP_CONFIG="$APP_CONFIG" go run ./cmd/server
