#!/bin/bash

# Script to create new migration files
# Usage: ./scripts/create-migration.sh migration_name

if [ -z "$1" ]; then
    echo "Usage: ./scripts/create-migration.sh migration_name"
    echo "Example: ./scripts/create-migration.sh add_users_avatar"
    exit 1
fi

MIGRATION_NAME=$1
MIGRATIONS_DIR="./migrations"

# Get the highest migration number
HIGHEST=$(ls -1 $MIGRATIONS_DIR | grep -o '^[0-9]\+' | sort -n | tail -1)

# Increment by 1
if [ -z "$HIGHEST" ]; then
    NEW_NUMBER="000001"
else
    NEW_NUMBER=$(printf "%06d" $((10#$HIGHEST + 1)))
fi

UP_FILE="$MIGRATIONS_DIR/${NEW_NUMBER}_${MIGRATION_NAME}.up.sql"
DOWN_FILE="$MIGRATIONS_DIR/${NEW_NUMBER}_${MIGRATION_NAME}.down.sql"

# Create migration files
echo "-- Add your UP migration here" > $UP_FILE
echo "-- Add your DOWN migration here" > $DOWN_FILE

echo "Created migration files:"
echo "  $UP_FILE"
echo "  $DOWN_FILE"
