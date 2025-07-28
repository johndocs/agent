#!/bin/bash

# This script runs the Go application twice and compares the output logs.
# It measures the consistency of the output by diffing the logs.

set -e

LOG_DIR="logs.directory_contents"
DIFF="$LOG_DIR/diff.txt"
rm -r "$LOG_DIR" 2>/dev/null || true
mkdir -p "$LOG_DIR"

# Function to create the command array.
# It takes one argument: the log number.
# It sets the global CMD array.
create_cmd_array() {
    local log_num=$1
    CMD=(
        go run . -d -o "$LOG_DIR" -l "$log_num"
        "What is this directory?"
        "Please examine the contents of all files"
        "Provide a short overview of this project for the README.md ## Summary"
    )
}

# The command and its arguments are defined in an array to improve readability and avoid long lines.
create_cmd_array 1
CMD1=("${CMD[@]}")
create_cmd_array 2
CMD2=("${CMD[@]}")

echo "Running first time... ${CMD1[@]}"
time "${CMD1[@]}"

echo "Running second time... ${CMD2[@]}"
time "${CMD2[@]}"

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
