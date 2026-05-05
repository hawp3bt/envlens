# envlens

> Diff and audit `.env` files across environments with secret masking.

---

## Installation

```bash
go install github.com/yourusername/envlens@latest
```

Or build from source:

```bash
git clone https://github.com/yourusername/envlens.git && cd envlens && go build -o envlens .
```

---

## Usage

Compare two `.env` files and mask sensitive values:

```bash
envlens diff .env.staging .env.production
```

**Example output:**

```
KEY                  STAGING              PRODUCTION
─────────────────────────────────────────────────────
DATABASE_URL         *** (masked)         *** (masked)
APP_PORT             8080                 9090          ← changed
LOG_LEVEL            debug                [missing]     ← missing
NEW_RELIC_KEY        [missing]            *** (masked)  ← added
```

Audit a single file for common issues (empty values, duplicates):

```bash
envlens audit .env.production
```

Show all keys without masking secrets (use with caution):

```bash
envlens diff .env .env.production --unmask
```

---

## Flags

| Flag        | Description                          |
|-------------|--------------------------------------|
| `--unmask`  | Reveal secret values in output       |
| `--json`    | Output results as JSON               |
| `--quiet`   | Only show differing keys             |

---

## Contributing

Pull requests are welcome. Please open an issue first to discuss any significant changes.

---

## License

[MIT](LICENSE) © 2024 yourusername