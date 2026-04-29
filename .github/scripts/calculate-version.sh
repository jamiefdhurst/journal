#!/bin/bash

# Get the latest tag, or default to 0.0.0 if no tags exist
LATEST_TAG=$(git describe --tags --abbrev=0 2>/dev/null || echo "0.0.0")
echo "Latest tag: $LATEST_TAG"

# Get commits since the latest tag
if [ "$LATEST_TAG" = "0.0.0" ]; then
  COMMITS=$(git log --pretty=format:"%s" --no-merges)
else
  COMMITS=$(git log ${LATEST_TAG}..HEAD --pretty=format:"%s" --no-merges)
fi

# Initialize bump type and releasable commits flag
BUMP_TYPE="patch"
HAS_RELEASABLE_COMMITS=false

# Analyze commits to determine bump type
while IFS= read -r commit; do
  # Skip empty lines
  [ -z "$commit" ] && continue

  # Convert to lowercase for case-insensitive matching
  commit_lower=$(echo "$commit" | tr '[:upper:]' '[:lower:]')

  # Skip commits that are only documentation, tests, chores, build, ci, style, or reverts
  if [[ "$commit_lower" =~ ^docs?(\(.*\))?:|^documentation:|^test(s)?(\(.*\))?:|^chore(\(.*\))?:|^build(\(.*\))?:|^ci(\(.*\))?:|^style(\(.*\))?:|^revert(\(.*\))?: ]] || \
     [[ "$commit_lower" =~ ^update.*documentation|^update.*readme|^update.*test ]]; then
    echo "Skipping non-release commit: $commit"
    continue
  fi

  # Mark that we have at least one releasable commit
  HAS_RELEASABLE_COMMITS=true

  # Check for major version indicators
  if [[ "$commit_lower" =~ ^breaking|^major|breaking\ change|breaking:|^[^:]+\!: ]]; then
    BUMP_TYPE="major"
    echo "Found breaking change: $commit"
    break
  fi

  # Check for minor version indicators (features, additions)
  if [[ "$commit_lower" =~ ^feat|^feature|^add|^minor ]]; then
    if [ "$BUMP_TYPE" != "major" ]; then
      BUMP_TYPE="minor"
      echo "Found feature/addition: $commit"
    fi
  fi

  # Everything else defaults to patch (fixes, updates, refactors, etc.)
done <<< "$COMMITS"

# If only non-releasable commits were found, don't create a release
if [ "$HAS_RELEASABLE_COMMITS" = false ]; then
  echo "Only non-releasable changes found (docs/test/chore/build/ci/style/revert) - no release needed"
  BUMP_TYPE="none"
fi

echo "Determined bump type: $BUMP_TYPE"

# Parse current version
if [ "$LATEST_TAG" = "0.0.0" ]; then
  MAJOR=0
  MINOR=0
  PATCH=0
else
  # Remove 'v' prefix if present
  VERSION=${LATEST_TAG#v}

  # Split version into components
  IFS='.' read -r MAJOR MINOR PATCH <<< "$VERSION"
fi

# Calculate new version based on bump type
case $BUMP_TYPE in
  major)
    MAJOR=$((MAJOR + 1))
    MINOR=0
    PATCH=0
    ;;
  minor)
    MINOR=$((MINOR + 1))
    PATCH=0
    ;;
  patch)
    PATCH=$((PATCH + 1))
    ;;
  none)
    # Keep version the same
    ;;
esac

NEW_VERSION="${MAJOR}.${MINOR}.${PATCH}"
echo "New version: $NEW_VERSION"

# Determine if we should create a release
if [ "$BUMP_TYPE" = "none" ] || [ "$LATEST_TAG" = "$NEW_VERSION" ] || [ "v$LATEST_TAG" = "$NEW_VERSION" ]; then
  SHOULD_RELEASE="false"
  echo "No release will be created (version unchanged)"
else
  SHOULD_RELEASE="true"
  echo "Release will be created"
fi

# Output to GitHub Actions
echo "version=$NEW_VERSION" >> $GITHUB_OUTPUT
echo "bump_type=$BUMP_TYPE" >> $GITHUB_OUTPUT
echo "previous_version=$LATEST_TAG" >> $GITHUB_OUTPUT
echo "should_release=$SHOULD_RELEASE" >> $GITHUB_OUTPUT
