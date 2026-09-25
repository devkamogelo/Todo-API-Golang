#!/usr/bin/env bash
set -e

# Load .env into the environment
set -a
source .env
set +a

command="$1"
name="$2"

case "$command" in
  up)
    migrate -path migrations -database "$DATABASE_URL" up
    ;;
  down)
    count="${name:-1}"
    read -p "Rolling back $count migration(s). Continue? [y/N] " confirm
    if [[ "$confirm" == "y" || "$confirm" == "Y" ]]; then
      migrate -path migrations -database "$DATABASE_URL" down "$count"
    fi
    ;;
  create)
    migrate create -ext sql -dir migrations -seq "$name"
    ;;
  force)
    migrate -path migrations -database "$DATABASE_URL" force "$name"
    ;;
  *)
    echo "Usage: $0 {up|down [n]|create <name>|force <version>}"
    exit 1
    ;;
esac