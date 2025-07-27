#!/bin/bash

# This script runs the Go application twice and compares the output logs.
# It measures the consistency of the output by diffing the logs.

set -e

LOG_DIR="logs"
LOG1="$LOG_DIR/log.1.txt"
LOG2="$LOG_DIR/log.2.txt"
LOGE1="$LOG_DIR/log.e.1.txt"
LOGE2="$LOG_DIR/log.e.2.txt"
DIFF="$LOG_DIR/diff.txt"
mkdir -p "$LOG_DIR"

# The command and its arguments are defined in an array to improve readability and avoid long lines.
CMD=(
    go run . -d
    "What is this directory?"
    "Please examine the contents of all files"
    "Provide a short overview of this project for the README.txt ## Summary"
)

echo "Running first time..."
time "${CMD[@]}" > "$LOG1" 2> "$LOGE1"

echo "Running second time..."
time "${CMD[@]}" > "$LOG2" 2> "$LOGE2"

echo
echo "---"
echo "Diffing stdout logs (log1.txt vs log2.txt):"
echo "---"
diff "$LOG1" "$LOG2" > "$DIFF" || true

echo
echo "---"
echo "Line counts for log files:"
echo "---"
# Find all files in the log directory and run wc on them.
# This avoids errors if the directory is empty.
find "$LOG_DIR" -type f -exec wc -l {} +
