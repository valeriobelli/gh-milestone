# gh milestone

A [gh](https://github.com/cli/cli) extension for managing Github Milestones

## Installation

```bash
gh extension install valeriobelli/gh-milestone
```

## Usage

Create a new milestone interactively

```bash
gh milestones create
```

or by using flags

```bash
gh milestones create --title v1.0.0 --description "This is a description" --due-date 2022-06-01
```

List milestones

```bash
gh milestone list
```

Edit a milestone

```bash
gh milestone edit <milestone number> --title "New title"
```

View a milestone

```bash
gh milestone view <milestone number>
```