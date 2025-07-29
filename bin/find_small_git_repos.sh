#!/bin/bash
# Script to find git repositories under ~ with fewer than 100 files

# Find all .git directories under home
find ~ -name .git -type d 2>/dev/null | while read gitdir; do
    # Get the parent directory (repository root)
    repo_path=$(dirname "$gitdir")
    
    # Count files in the repository
    file_count=$(find "$repo_path" -type f | wc -l)
    
    # Check if file count is less than or equal to 100
    if [ "$file_count" -le 100 ]; then
        echo "$repo_path"
    fi
done