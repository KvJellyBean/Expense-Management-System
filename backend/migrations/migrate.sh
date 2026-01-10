#!/bin/bash

# Database connection info
DB_HOST="${DB_HOST:-postgres}"
DB_PORT="${DB_PORT:-5432}"
DB_USER="${DB_USER:-expense_user}"
DB_PASSWORD="${DB_PASSWORD:-expense_pass}"
DB_NAME="${DB_NAME:-expense_db}"

export PGPASSWORD="$DB_PASSWORD"

echo "Waiting for database to be ready..."
until psql -h "$DB_HOST" -U "$DB_USER" -d "$DB_NAME" -c '\q' 2>/dev/null; do
  sleep 1
done

echo "Running migrations..."

for migration in /migrations/*.up.sql; do
  if [ -f "$migration" ]; then
    echo "Applying $(basename $migration)..."
    psql -h "$DB_HOST" -U "$DB_USER" -d "$DB_NAME" -f "$migration"
  fi
done

echo "Migrations completed!"
