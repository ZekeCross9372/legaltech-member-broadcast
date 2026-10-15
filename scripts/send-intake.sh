#!/bin/sh
set -eu

curl --fail-with-body --request POST "${LEGAL_BROADCAST_URL:-http://localhost:8080}/announcements" \
  --header 'Content-Type: application/json' \
  --data '{"matter_id":"MAT-42","stage":"matter_intake","summary":"New conflicts review","account_id":"legal-ops"}'
