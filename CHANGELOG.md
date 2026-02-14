# gh-milestone

## 2.2.0

### Minor Changes

- Add `--all` flag to `list` subcommand to retrieve all milestones via cursor-based pagination, bypassing the 100-item limit. Resolves [#27](https://github.com/valeriobelli/gh-milestone/issues/27).
- Add `openIssues` and `closedIssues` JSON fields to `list` subcommand output. Resolves [#31](https://github.com/valeriobelli/gh-milestone/issues/31).
- Add `TITLE` and `ISSUES` as `--orderBy.field` options with client-side sorting. Resolves [#13](https://github.com/valeriobelli/gh-milestone/issues/13).
- Support `gh config set editor` for editor resolution in interactive mode.

### Patch Changes

- Remove debug output that corrupted JSON when using `--repo`. Resolves [#30](https://github.com/valeriobelli/gh-milestone/issues/30).
- Remove hardcoded `--wait` flag from editor invocation. Resolves [#32](https://github.com/valeriobelli/gh-milestone/issues/32).
- Upgrade dependencies to resolve all known vulnerabilities (`go-gh` v2, `go-github` v68, inline `surveyext` editor package).

## 2.1.0

### Minor Changes

- Add --repo/-R option to every command. Resolves [#15](https://github.com/valeriobelli/gh-milestone/issues/15)

## 2.0.1

### Patch Changes

- Consider timezone when setting the Due date. Resolves [#19](https://github.com/valeriobelli/gh-milestone/issues/19).
- Remove unwanted white foreground color when printing a Milestone's status. Resolves [#20](https://github.com/valeriobelli/gh-milestone/issues/20).

## 2.0.0

### Major Changes

- Remove the `--output=json` flag from the `list` subcommand. The specific JSON output is now exposed via the new `--json` flag which acts as the other gh/cli subcommands.
- The JSON output returned by `list` subcommand now contains camel case attributes.

### Patch Changes

- Resolve `jq` bugged syntax

## 1.1.1

### Patch Changes

- Strings with mixed cases are now recognized and handled in `edit` subcommand.

## 1.1.0

### Minor Changes

- Enhance the interface of commands. Strict values are now validated and better error messages are emitted.

## 1.0.0

### Major Changes

- Rename the `states` option of `list` subcommand to `state`. It now takes `all`, `closed` and `open` as possible values.
  
  ```bash
  gh milestone list --state all
  gh milestone list --state closed
  gh milestone list --state open
  ```

### Minor Changes

- Tidy up the `list` subcommand's helper.
