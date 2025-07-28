#!/bin/bash

# This script runs the Go application twice and compares the output logs.
# It measures the consistency of the output by diffing the logs.

set -e

LOG_DIR="logs.search"
DIFF="$LOG_DIR/diff.txt"
rm -r "$LOG_DIR" 2>/dev/null || true
mkdir -p "$LOG_DIR"

# The command and its arguments are defined in an array to improve readability and avoid long lines.
# CMD=(
#     go run . -d
#     "Find all directories under ~ that are git repositories and print their paths without the trailing '/.git'. Exclude any directories that contain more than 100 files. They are likely to be clones of large libraries."
#     "Tell me about each git repositories you found. What programs and/or tools does it implement?  For each program, provide a short description of its purpose and functionality. Then suggest how it could be used. Save this information in a directory 'repositories'. For each repository, create a json file called 'repositories/repo.<repo_name>.json' so that you don't need to use up your token quota. For each repository, also provide a list of the files in the repository and their sizes in bytes. Save this information in a json file called 'repositories/files.<repo_name>.json'"
#     "Create a simple Go program that does a plain text search over all the repositories found. It should respond to queries about the intent of the code, the purpose of the code, and how to use the code. It should also be able to answer questions about the code's functionality and implementation details."
# )

CMD=(
    go run . -d -o "$LOG_DIR" -l 1
    "Create a simple Go program 'bin/repo_search.go' that does a plain text search over all the repositories in the './repositories' directory. It should respond to queries about the intent of the code, the purpose of the code, and how to use the code. It should also be able to answer questions about the code's functionality and implementation details. If 'bin/repo_search.go' already exists, you started it in a previous run, so just continue where you left off. If it doesn't exist, create it from scratch. The program should be able to handle large repositories and should be efficient in its search algorithm."
)

echo "Building search app..."
time "${CMD[@]}"


echo
echo "---"
echo "Line counts for log files:"
echo "---"
# Find all files in the log directory and run wc on them.
# This avoids errors if the directory is empty.
find "$LOG_DIR" -type f -exec wc -l {} +
