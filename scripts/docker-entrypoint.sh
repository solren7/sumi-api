#!/bin/sh
set -e

if [ "${AUTO_MIGRATE}" = "true" ]; then
  echo "Running auto migration before startup..."
  /app/server migrate
fi

exec /app/server api
