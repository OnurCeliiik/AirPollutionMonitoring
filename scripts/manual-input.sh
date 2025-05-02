#!/bin/bash

# Manual input script for testing the air quality monitoring platform
# Usage: ./manual-input.sh <latitude> <longitude> <parameter> <value>
# Example: ./manual-input.sh 41.015 28.979 "PM2.5" 90.0

# Default values
API_HOST="${API_HOST:-http://localhost:8080}"
CURRENT_TIME=$(date -u +"%Y-%m-%dT%H:%M:%SZ")

# Check if required parameters are provided
if [ "$#" -lt 4 ]; then
    echo "Usage: $0 <latitude> <longitude> <parameter> <value>"
    echo "Example: $0 41.015 28.979 \"PM2.5\" 90.0"
    exit 1
fi

# Read parameters
LATITUDE=$1
LONGITUDE=$2
PARAMETER=$3
VALUE=$4
TIMESTAMP="${5:-$CURRENT_TIME}"

# Construct JSON payload
JSON_PAYLOAD=$(cat << EOF
{
  "latitude": $LATITUDE,
  "longitude": $LONGITUDE,
  "parameter": "$PARAMETER",
  "value": $VALUE,
  "timestamp": "$TIMESTAMP"
}
EOF
)

# Send POST request
echo "Sending data: $JSON_PAYLOAD"
curl -X POST \
  -H "Content-Type: application/json" \
  -d "$JSON_PAYLOAD" \
  $API_HOST/api/data

echo "" 