#!/usr/bin/env bash

# wait-for-it.sh -- host:port

set -e

host="$1"
shift

until nc -z "$host" 3306; do
  >&2 echo "MySQL ainda não está pronto - aguardando..."
  sleep 1
done

>&2 echo "MySQL está pronto - executando comando"
exec "$@" 