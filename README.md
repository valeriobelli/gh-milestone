# gh milestone

A [gh](https://github.com/cli/cli) extension for managing GitHub Milestones.

## Installation

```bash
gh extension install valeriobelli/gh-milestone
```

## Usage

By default, this extension uses the GitHub's Access Token of the current user for the host `github.com`.

You can override the authentication of the current registered gh user by defining `GITHUB_TOKEN` in your environment variables. You can also override the host where your user is authenticated by setting `GITHUB_MILESTONE_HOST` in your environment variables.

### Editor configuration

When creating a milestone interactively, the extension opens an editor for the description field. The editor is resolved in the following priority order:

1. `GH_EDITOR` environment variable
2. `gh config set editor` value (e.g. `gh config set editor "code --wait"`)
3. `GIT_EDITOR` environment variable
4. `VISUAL` environment variable
5. `EDITOR` environment variable
6. `nano` (or `notepad` on Windows)

> **Tip:** If you use a GUI editor like VS Code, make sure to include the `--wait` flag so the extension waits for you to close the file (e.g. `gh config set editor "code --wait"`).

### Create a new milestone

```bash
# Interactive mode
gh milestone create

# Flags mode
gh milestone create --title v1.0.0 --description "This is a description" --due-date 2022-06-01
```

### List milestones

```bash
# Extended command
gh milestone list

# Alias
gh milestone ls

# List closed Milestones
gh milestone list --state closed

# List milestones of specific repo
gh milestone list --repo valeriobelli/foo-bar

# Search by a pattern
gh milestone list --query "Foo bar"

# Get first ten milestones
gh milestone list --first 10

# Get all milestones (bypasses the 100 item limit)
gh milestone list --all

# Sort by title
gh milestone list --orderBy.field title

# Sort by issue count (descending)
gh milestone list --orderBy.field issues --orderBy.direction desc

# Print milestones as JSON
gh milestone list --json id
gh milestone list --json id,progressPercentage --json number

# Print open and closed issue counts as JSON
gh milestone list --json openIssues,closedIssues

# Access Milestone attributes via jq
gh milestone list --json id,progressPercentage --json number --jq ".[0].id"
```

#### Sorting

Milestones can be sorted using `--orderBy.field` and `--orderBy.direction`:

| Field        | Description                        |
| ------------ | ---------------------------------- |
| `created_at` | Sort by creation date              |
| `due_date`   | Sort by due date                   |
| `issues`     | Sort by total issue count          |
| `number`     | Sort by milestone number (default) |
| `title`      | Sort by title (alphabetical)       |
| `updated_at` | Sort by last update date           |

Direction can be `asc` (default) or `desc`.

#### Pagination

By default, up to 100 milestones are returned. Use `--first N` to change the limit, or `--all` to retrieve every milestone regardless of count.

```bash
gh milestone list --first 20
gh milestone list --all
```

### Edit a milestone

```bash
gh milestone edit <milestone number> --title "New title"
gh milestone edit <milestone number> --title "New title" --repo valeriobelli/foo-bar
```

### View a milestone

```bash
gh milestone view <milestone number>

gh milestone view <milestone number> --repo valeriobelli/foo-bar
```

### Delete milestone

```bash
# Interactive mode
gh milestone delete <milestone number>

# Automatic
gh milestone delete <milestone number> --confirm
gh milestone delete <milestone number> --confirm --repo valeriobelli/foo-bar
```
