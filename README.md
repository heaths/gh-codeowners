# CODEOWNERS GitHub CLI extension

[![releases](https://img.shields.io/github/v/release/heaths/gh-codeowners.svg?logo=github)](https://github.com/heaths/gh-codeowners/releases/latest)
[![ci](https://github.com/heaths/gh-codeowners/actions/workflows/ci.yml/badge.svg?event=push)](https://github.com/heaths/gh-codeowners/actions/workflows/ci.yml)

Lint your CODEOWNERS file.

## Usage

Render unknown owners red based on the current branch's CODEOWNERS errors reported by GitHub:

```bash
gh codeowners lint
```

You can also get a sorted list of unknown owners:

```bash
gh codeowners lint --unknown-owners
```

## License

Licensed under the [MIT](LICENSE.txt) license.
