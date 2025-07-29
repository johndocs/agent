#!/bin/bash
# Script to analyze git repositories and save information to JSON files

# Get the list of repositories
repos=$(./bin/find_small_git_repos.sh)

for repo_path in $repos; do
    # Extract repo name from path
    repo_name=$(basename "$repo_path")
    echo "Analyzing repository: $repo_name"
    
    # Create files JSON
    echo "{" > "repositories/files.$repo_name.json"
    echo "  \"repository\": \"$repo_name\"," >> "repositories/files.$repo_name.json"
    echo "  \"path\": \"$repo_path\"," >> "repositories/files.$repo_name.json"
    echo "  \"files\": [" >> "repositories/files.$repo_name.json"
    
    first_file=true
    while IFS= read -r file_info; do
        if [ "$first_file" = true ]; then
            first_file=false
        else
            echo "," >> "repositories/files.$repo_name.json"
        fi
        
        file_path=$(echo "$file_info" | awk '{print $2}')
        file_size=$(echo "$file_info" | awk '{print $1}')
        
        # Remove repo path prefix to get relative path
        relative_path=${file_path#$repo_path/}
        
        echo "    {" >> "repositories/files.$repo_name.json"
        echo "      \"path\": \"$relative_path\"," >> "repositories/files.$repo_name.json"
        echo "      \"size\": $file_size" >> "repositories/files.$repo_name.json"
        echo -n "    }" >> "repositories/files.$repo_name.json"
    done < <(find "$repo_path" -type f -not -path "*/\.*" -exec du -b {} \;)
    
    echo "" >> "repositories/files.$repo_name.json"
    echo "  ]" >> "repositories/files.$repo_name.json"
    echo "}" >> "repositories/files.$repo_name.json"
    
    # Basic repo info for now - we'll analyze contents later
    echo "{" > "repositories/repo.$repo_name.json"
    echo "  \"name\": \"$repo_name\"," >> "repositories/repo.$repo_name.json"
    echo "  \"path\": \"$repo_path\"," >> "repositories/repo.$repo_name.json"
    echo "  \"analysis\": {" >> "repositories/repo.$repo_name.json"
    echo "    \"description\": \"\"," >> "repositories/repo.$repo_name.json"
    echo "    \"programs\": []," >> "repositories/repo.$repo_name.json"
    echo "    \"suggestion\": \"\"" >> "repositories/repo.$repo_name.json"
    echo "  }" >> "repositories/repo.$repo_name.json"
    echo "}" >> "repositories/repo.$repo_name.json"
done