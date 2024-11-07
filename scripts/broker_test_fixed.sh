#!/bin/bash

# Enable error handling
set -e

TICK_DIR_NAME="./log/tickLogger/"
HOP_DIR_NAME="./log/hopLogger/"
MAX_INDEX=0
RUN="./cmd/broker/run/main.go"
BROKER_START=50001
BROKER_END=50005

# Check if directories exist
if [ ! -d "$TICK_DIR_NAME" ]; then
    echo "Directory $TICK_DIR_NAME not found."
    exit 1
fi

if [ ! -d "$HOP_DIR_NAME" ]; then
    echo "Directory $HOP_DIR_NAME not found."
    exit 1
fi

# Find the highest index
while [ -d "${TICK_DIR_NAME}${MAX_INDEX}" ]; do
    MAX_INDEX=$((MAX_INDEX + 1))
    echo $MAX_INDEX
done

# Function to start broker service in a new terminal window on macOS
start_broker_service() {
    local port=$1
    osascript <<EOF
tell application "Terminal"
    activate
    do script "cd $PWD && go run $RUN --port $port --dir_index $MAX_INDEX"
end tell
EOF
}

# Start broker services
for ((P=BROKER_START; P<=BROKER_END; P++)); do
    echo "Starting Broker Service on Port $P"
    start_broker_service $P
done

wait
