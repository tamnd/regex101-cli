# regex101

Browse the [regex101.com](https://regex101.com) shared pattern library from the command line.

`regex101` is a single pure-Go binary. No API key required.

## Install

```bash
go install github.com/tamnd/regex101-cli/cmd/regex101@latest
```

Or grab a prebuilt binary from the [releases](https://github.com/tamnd/regex101-cli/releases), or run the container image:

```bash
docker run --rm ghcr.io/tamnd/regex101:latest --help
```

## Usage

```bash
# List top regex patterns by upvotes (default 20)
regex101 list

# Filter by language flavor
regex101 list --flavor python
regex101 list --flavor javascript
regex101 list --flavor pcre

# Search patterns by keyword
regex101 search "email"
regex101 search "url validation" --flavor javascript

# Output formats
regex101 list -o json
regex101 list -o csv
regex101 search "phone" -o table
```

## Commands

| Command | Description |
|---------|-------------|
| `list` | List top regex patterns by upvotes (optional `--flavor` filter) |
| `search <query>` | Search patterns by keyword (optional `--flavor` filter) |
| `version` | Show version information |

## Flavors

`javascript`, `python`, `pcre`, `pcre2`, `php`, `golang`, `java`, `ruby`, `rust`, `csharp`

## Global flags

```
-o, --output string    output format: table|json|jsonl|csv|tsv|url|raw (default "auto")
-n, --limit int        limit number of records (0 = command default)
    --fields strings   comma-separated columns to include
    --no-header        omit header row
    --template string  Go text/template per record
    --timeout duration per-request timeout (default 30s)
    --delay duration   minimum spacing between requests
```
