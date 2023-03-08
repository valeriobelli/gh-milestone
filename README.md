# gh milestone

A [gh](https://github.com/cli/cli) extension for managing Github Milestones

## Installation

```bash
gh extension install valeriobelli/gh-milestone
```

## Usage

By default, this extension uses the GitHub's Access Token of the current user for the host `github.com`.

You override the authentication of the current registered gh user by defining `GITHUB_TOKEN` in your environment variables. You can also override the host where your user is authenticated by setting `GITHUB_MILESTONE_HOST` still in your environment variables. By defining this latter, you change the behaviour of `gh config get --host <host> oauth_token`, which is the command this extension relies upon.

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

# Print milestones as JSON
gh milestone list --json id
gh milestone list --json id,progressPercentage --json number

# Access Milestone attributes via jq
gh milestone list --json id,progressPercentage --json number --jq ".[0].id"
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
