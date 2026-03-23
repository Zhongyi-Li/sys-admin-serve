#!/usr/bin/env sh
set -eu

ROOT_DIR=$(CDPATH= cd -- "$(dirname "$0")/.." && pwd)
MIGRATIONS_DIR="$ROOT_DIR/migrations"
NAME=${1:-}

if [ -z "$NAME" ]; then
	echo "usage: sh ./scripts/new-migration.sh <name>" >&2
	exit 1
fi

timestamp=$(date +%Y%m%d%H%M%S)
up_file="$MIGRATIONS_DIR/${timestamp}_${NAME}.up.sql"
down_file="$MIGRATIONS_DIR/${timestamp}_${NAME}.down.sql"

touch "$up_file" "$down_file"

echo "created: $up_file"
echo "created: $down_file"