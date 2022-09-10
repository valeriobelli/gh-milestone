# gh-milestone

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
