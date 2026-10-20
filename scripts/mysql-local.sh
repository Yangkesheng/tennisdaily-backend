#!/usr/bin/env bash
# Use the Docker MySQL client and Unix socket; no host port mapping is required.
set -euo pipefail

if [[ "${1:-}" == "--help" || "${1:-}" == "-h" ]]; then
  cat <<'HELP'
Usage:
  ./scripts/mysql-local.sh
  ./scripts/mysql-local.sh -e 'SHOW TABLES;'
  ./scripts/mysql-local.sh < migrations/example.sql

Connects inside tennisdaily-mysql using a Unix socket, user root, DB tennis_diary.
Overrides: LOCAL_MYSQL_CONTAINER, LOCAL_MYSQL_USER, LOCAL_MYSQL_DATABASE.
Automatically uses container credentials; an existing MYSQL_PWD takes priority.
Requires Docker and a running MySQL container. No host mysql client is needed.
HELP
  exit 0
fi

if ! command -v docker >/dev/null 2>&1; then
  echo 'docker not found. Start Docker Desktop and check PATH.' >&2
  exit 1
fi

exec_args=(-i)
if [[ -t 0 && -t 1 ]]; then
  exec_args+=(-t)
fi
if [[ "${MYSQL_PWD+x}" == x ]]; then
  exec_args+=(-e MYSQL_PWD)
fi

exec docker exec "${exec_args[@]}" "${LOCAL_MYSQL_CONTAINER:-tennisdaily-mysql}" sh -c '
  set -eu
  if [ "${MYSQL_PWD+x}" != x ]; then
    if [ "$1" = root ]; then
      if [ -n "${MYSQL_ROOT_PASSWORD_FILE:-}" ]; then
        MYSQL_PWD=$(cat "$MYSQL_ROOT_PASSWORD_FILE")
      elif [ "${MYSQL_ROOT_PASSWORD+x}" = x ]; then
        MYSQL_PWD=$MYSQL_ROOT_PASSWORD
      elif [ "${MYSQL_ALLOW_EMPTY_PASSWORD:-}" = yes ]; then
        MYSQL_PWD=
      else
        echo "Container has no configured root password." >&2
        exit 1
      fi
    elif [ "$1" = "${MYSQL_USER:-}" ]; then
      if [ -n "${MYSQL_PASSWORD_FILE:-}" ]; then
        MYSQL_PWD=$(cat "$MYSQL_PASSWORD_FILE")
      else
        MYSQL_PWD=${MYSQL_PASSWORD:-}
      fi
    else
      echo "Requested user does not match a configured container user." >&2
      exit 1
    fi
  fi
  export MYSQL_PWD
  db_user=$1
  db_name=$2
  shift 2
  exec mysql --no-defaults --protocol=SOCKET --user="$db_user" --database="$db_name" --default-character-set=utf8mb4 "$@"
' sh "${LOCAL_MYSQL_USER:-root}" "${LOCAL_MYSQL_DATABASE:-tennis_diary}" "$@"
