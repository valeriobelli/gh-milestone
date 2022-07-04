# gh milestone

A [gh](https://github.com/cli/cli) extension for managing Github Milestones

## Installation

```bash
gh extension install valeriobelli/gh-milestone
```

## Usage

By default, this extension uses the GitHub's Access Token of the current user for the host `github.com`.
You can override this behaviour by setting either `GITHUB_TOKEN` or `GITHUB_MILESTONE_HOST`.

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
```

### Edit a milestone

```bash
gh milestone edit <milestone number> --title "New title"
```

### View a milestone

```bash
gh milestone view <milestone number>
```

### Delete milestone

```bash
# Interactive mode
gh milestone delete <milestone number>

# Automatic
gh milestone delete <milestone number> --confirm
```