#!/bin/bash

# Extract release notes for a given version from CHANGELOG.md

CHANGELOG="CHANGELOG.md"
OUTPUT="release_notes.md"

# Function to extract notes for a specific version
extract_notes() {
    local version=$1
    local in_section=false
    local notes=""

    while IFS= read -r line; do
        # Check if we've found the version header
        if [[ $line =~ ^##[[:space:]]\[$version\] ]]; then
            in_section=true
            continue
        fi

        # Check if we've hit the next version header
        if [[ $in_section == true ]] && [[ $line =~ ^##[[:space:]]\[ ]]; then
            break
        fi

        # Collect lines if we're in the section
        if [[ $in_section == true ]]; then
            notes+="$line"$'\r\n'
        fi
    done < "$CHANGELOG"

    # Trim trailing newlines
    echo "$notes" | sed -e :a -e '/^\n*$/{$d;N;ba' -e '}'
}

echo $(extract_notes $1) > $OUTPUT
