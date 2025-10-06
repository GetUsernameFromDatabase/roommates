#!/bin/bash
set -e

SCRIPT_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )
cd $SCRIPT_DIR

# ------------------------------------------------------------------------------
FOLDER="$SCRIPT_DIR/mock-sql"
if [ "$#" -lt 3 ]; then
  echo "Usage: $0 <db_user> <db_password> <db_name> <sql_folder> [host] [port]"
  exit 1
fi

DB_USER=$1
DB_PASSWORD=$2
DB_NAME=$3
DB_HOST=${4:-localhost}
DB_PORT=${5:-5432}

export PGPASSWORD="$DB_PASSWORD"

# ------------------------------------------------------------------------------

echo "🔗 Will be running plsql using $DB_USER@$DB_HOST:$DB_PORT/$DB_NAME"
echo "📁 Looking for sql files under: $FOLDER"

for sql_file in "$FOLDER"/*.sql; do
  [ -e "$sql_file" ] || { echo "⚠️ No .sql files found in $FOLDER"; break; }

  echo "📄 Executing: $sql_file"
  psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" -f "$sql_file"
done
echo "✅ All SQL files executed."
