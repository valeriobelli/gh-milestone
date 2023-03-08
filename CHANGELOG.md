# gh-milestone

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
