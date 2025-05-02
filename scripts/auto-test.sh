#!/bin/bash

# Automated test script for the air quality monitoring platform
# Usage: ./auto-test.sh [--duration=<seconds>] [--rate=<requests_per_minute>] [--anomaly-chance=<percent>]
# Example: ./auto-test.sh --duration=60 --rate=10 --anomaly-chance=20

# Default values
API_HOST="${API_HOST:-http://localhost:8080}"
DURATION=60  # seconds
RATE=10      # requests per minute
ANOMALY_CHANCE=20  # percent chance of generating an anomaly

# Parse named parameters
for i in "$@"; do
    case $i in
        --duration=*)
        DURATION="${i#*=}"
        shift
        ;;
        --rate=*)
        RATE="${i#*=}"
        shift
        ;;
        --anomaly-chance=*)
        ANOMALY_CHANCE="${i#*=}"
        shift
        ;;
        *)
        # unknown option
        ;;
    esac
done

echo "Starting automated test with:"
echo "  Duration: $DURATION seconds"
echo "  Rate: $RATE requests per minute"
echo "  Anomaly chance: $ANOMALY_CHANCE%"

# Calculate delay between requests
DELAY=$(bc <<< "scale=3; 60 / $RATE")

# Array of parameters to monitor
PARAMETERS=("PM2.5" "PM10" "NO2" "O3")

# Normal value ranges for each parameter (min max)
PM25_RANGE=(5 30)
PM10_RANGE=(10 50)
NO2_RANGE=(5 40)
O3_RANGE=(30 100)

# Anomaly value ranges for each parameter (min max)
PM25_ANOMALY=(50 200)
PM10_ANOMALY=(70 250)
NO2_ANOMALY=(60 180)
O3_ANOMALY=(150 300)

# Generate a random location around a central point
# Args: latitude longitude radius_km
function random_location() {
    local lat=$1
    local lon=$2
    local radius=$3
    
    # Earth's radius in km
    local earth_radius=6371
    
    # Convert radius from km to degrees (approximate)
    local radius_deg=$(bc <<< "scale=6; $radius / $earth_radius * (180 / 3.14159)")
    
    # Random angle
    local angle=$(bc <<< "scale=6; $RANDOM / 32768 * 2 * 3.14159")
    
    # Random distance within radius
    local distance=$(bc <<< "scale=6; sqrt($RANDOM / 32768) * $radius_deg")
    
    # Calculate new point
    local new_lat=$(bc <<< "scale=6; $lat + $distance * cos($angle)")
    local new_lon=$(bc <<< "scale=6; $lon + $distance * sin($angle)")
    
    echo "$new_lat $new_lon"
}

# Generate a random value between min and max
# Args: min max
function random_value() {
    local min=$1
    local max=$2
    local range=$(bc <<< "$max - $min")
    local value=$(bc <<< "scale=1; $min + $RANDOM / 32768 * $range")
    echo $value
}

# Reference location (city center coordinates)
CENTER_LAT=41.015
CENTER_LON=28.979

# Start time
START_TIME=$(date +%s)
END_TIME=$((START_TIME + DURATION))

# Run until duration is reached
while [ $(date +%s) -lt $END_TIME ]; do
    # Select a random parameter
    PARAM_INDEX=$((RANDOM % ${#PARAMETERS[@]}))
    PARAMETER=${PARAMETERS[$PARAM_INDEX]}
    
    # Generate random location within 10km of center
    LOCATION=$(random_location $CENTER_LAT $CENTER_LON 10)
    LAT=$(echo $LOCATION | cut -d' ' -f1)
    LON=$(echo $LOCATION | cut -d' ' -f2)
    
    # Decide if this should be an anomaly
    IS_ANOMALY=$((RANDOM % 100 < ANOMALY_CHANCE))
    
    # Generate value based on parameter and anomaly status
    case $PARAMETER in
        "PM2.5")
            if [ $IS_ANOMALY -eq 1 ]; then
                VALUE=$(random_value ${PM25_ANOMALY[0]} ${PM25_ANOMALY[1]})
            else
                VALUE=$(random_value ${PM25_RANGE[0]} ${PM25_RANGE[1]})
            fi
            ;;
        "PM10")
            if [ $IS_ANOMALY -eq 1 ]; then
                VALUE=$(random_value ${PM10_ANOMALY[0]} ${PM10_ANOMALY[1]})
            else
                VALUE=$(random_value ${PM10_RANGE[0]} ${PM10_RANGE[1]})
            fi
            ;;
        "NO2")
            if [ $IS_ANOMALY -eq 1 ]; then
                VALUE=$(random_value ${NO2_ANOMALY[0]} ${NO2_ANOMALY[1]})
            else
                VALUE=$(random_value ${NO2_RANGE[0]} ${NO2_RANGE[1]})
            fi
            ;;
        "O3")
            if [ $IS_ANOMALY -eq 1 ]; then
                VALUE=$(random_value ${O3_ANOMALY[0]} ${O3_ANOMALY[1]})
            else
                VALUE=$(random_value ${O3_RANGE[0]} ${O3_RANGE[1]})
            fi
            ;;
    esac
    
    # Current time in ISO format
    TIMESTAMP=$(date -u +"%Y-%m-%dT%H:%M:%SZ")
    
    # Construct JSON payload
    JSON_PAYLOAD=$(cat << EOF
{
  "latitude": $LAT,
  "longitude": $LON,
  "parameter": "$PARAMETER",
  "value": $VALUE,
  "timestamp": "$TIMESTAMP"
}
EOF
    )
    
    # Send POST request
    if [ $IS_ANOMALY -eq 1 ]; then
        echo "Sending ANOMALY data: $JSON_PAYLOAD"
    else
        echo "Sending normal data: $JSON_PAYLOAD"
    fi
    
    curl -s -X POST \
      -H "Content-Type: application/json" \
      -d "$JSON_PAYLOAD" \
      $API_HOST/api/data > /dev/null
    
    # Wait for next request
    sleep $DELAY
done

echo "Test completed after $DURATION seconds." 