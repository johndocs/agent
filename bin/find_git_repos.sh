#!/bin/bash

# Script to find all git repositories under the home directory
# Excludes repositories with more than 100 files

echo "Finding git repositories under home directory..."
echo "Excluding repositories with more than 100 files..."

# Find all .git directories
find ~ -name .git -type d 2>/dev/null | while read git_dir; do
    # Get the parent directory (actual repo directory)
    repo_dir=$(dirname "$git_dir")
    
    # Count files in the repository
    file_count=$(find "$repo_dir" -type f | wc -l)
    
    # Check if file count is less than or equal to 100
    if [ "$file_count" -le 100 ]; then
        echo "$repo_dir"
    fi
done

echo "Search complete."