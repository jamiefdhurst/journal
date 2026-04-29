# Commit Message Convention

This project uses a commit message convention to automatically determine version bumps and generate changelogs.

## Format

All commit messages must follow this format:

```
<type>[optional scope]: <description>
```

## Commit Types

### Version Bumps

These commit types will trigger a new release:

| Type                        | Version Bump      | Description                    | Example                           |
| --------------------------- | ----------------- | ------------------------------ | --------------------------------- |
| `feat:`, `feature:`, `add:` | **Minor** (0.X.0) | New features or additions      | `feat: Add user authentication`   |
| `fix:`, `bug:`              | **Patch** (0.0.X) | Bug fixes                      | `fix: Resolve login issue`        |
| `update:`, `refactor:`      | **Patch** (0.0.X) | Code improvements              | `refactor: Simplify API client`   |
| `improve:`, `perf:`         | **Patch** (0.0.X) | Performance improvements       | `perf: Optimize database queries` |
| `breaking:`, `major:`       | **Major** (X.0.0) | Breaking changes               | `breaking: Remove deprecated API` |
| `<type>!:`                  | **Major** (X.0.0) | Breaking change (conventional) | `feat!: Redesign authentication`  |

### No Version Bump

These commit types will NOT trigger a release:

| Type              | Description             | Example                        |
| ----------------- | ----------------------- | ------------------------------ |
| `docs:`, `doc:`   | Documentation changes   | `docs: Update README`          |
| `test:`, `tests:` | Test changes only       | `test: Add integration tests`  |
| `chore:`          | Maintenance tasks       | `chore: Update dependencies`   |
| `build:`          | Build system changes    | `build: Update webpack config` |
| `ci:`             | CI/CD changes           | `ci: Add new workflow step`    |
| `style:`          | Code style/formatting   | `style: Run prettier`          |
| `revert:`         | Revert previous commits | `revert: Revert feature X`     |

## Scopes (Optional)

You can add a scope to provide more context:

```
feat(auth): Add OAuth support
fix(api): Resolve timeout issue
chore(deps): Update packages
```

## Examples

### Good Commit Messages

```
feat: Add dark mode toggle
fix(ui): Resolve button alignment issue
docs: Update installation instructions
chore(deps): Bump lodash from 4.17.19 to 4.17.21
refactor: Extract validation logic
perf: Optimize image loading
breaking: Remove support for Node 12
build: Update webpack to v5
ci: Add code coverage reporting
style: Apply prettier formatting
revert: Revert "Add experimental feature"
```

### Bad Commit Messages (Will be Rejected)

```
Update files
WIP
quick fix
Fixed bug
Added feature
```

## Automatic Validation

Commit messages are validated in two places:

1. **Locally (git hook)**: When you commit, a git hook validates your message. If invalid, you'll see a helpful error message with examples.
2. **In CI (GitHub Actions)**: All commits in a pull request are validated using the same hook. The PR will fail if any commits don't follow the convention.

This ensures everyone follows the same standards, whether they have git hooks installed or not.

## Version Calculation

- Pull requests show a comment indicating what version bump will occur
- When merged to main, the version is automatically calculated and released
- Only commits with releasable types trigger new versions

## Special Cases

- Merge commits are automatically skipped from validation
- `[skip ci]` commits (auto-generated) are skipped from validation
- Only non-releasable changes (docs/test/chore/build/ci/style/revert) won't trigger a release
- If a PR contains both releasable and non-releasable commits, a release will be triggered
- Dependabot is configured to use `chore(deps):` format automatically

