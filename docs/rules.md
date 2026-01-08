# Writing Rules

Rules describe **what production risk looks like**.  
They do not perform detection — they evaluate **signals** produced by detectors.

A rule answers one question:

> “If this signal pattern is true, what operational risk does it represent?”

---

## Rule structure

Example:

```yaml
id: secrets-management
severity: high
category: security
title: Secrets likely stored as environment variables

description: >
  Secrets appear to be handled via environment variables without a
  dedicated secrets management solution.

why_it_matters:
  - Environment variables are often logged, dumped, or exposed by mistake.
  - Rotating env-based secrets usually requires redeployments.
  - Access control and auditing are typically missing.

detect:
  any_of:
    - file_exists: .env
    - code_contains: process.env
    - code_contains: os.environ
    - signal_equals:
        secrets_provider_detected: false

confidence: high

```

## Required fields

| Field | Description   |
|-----|----|
| `id` | Unique identifier of the rule  |
| `severity`      | `high`, `medium`, or `low`     |
| `category`    | Logical grouping (e.g. `deployment`, `security`) |
| `title`     | Human-readable summary  |
| `Description`     | What was detected  |
| `why_it_matters`     | Human-readable summary  |
| `detect`       | Logical conditions  |

## Detect conditions

Rules use logical operators:

```yaml
detect:
  any_of:   # OR
  all_of:   # AND
  none_of:  # NOT
```

## Supported conditions

| Condition | Meaning   |
|-----|----|
| `file_exists` | A file exists in the repo  |
| `code_contains`    | A string appears in scanned code |
| `signal_equals`       | A detected signal has a specific value  |

## Example conditions

File check

```yaml
file_exists: .env
```

Code pattern

```yaml
code_contains: process.env
```

Signal check

```yaml
signal_equals:
  secrets_provider_detected: true
```

## Adding a new rule

1. Create a new YAML file in **rules** foleder
2. Give it a unique id
3. Reference existing signals
4. Run:

```bash
pr scan .
```

5. Review the report

## Rule design principles

* Prefer false negatives over false positives
* Be risk-oriented, not style-oriented
* Always explain why it matters
* Avoid vendor-specific assumptions
* Write as if teaching a junior engineer

Rules are not policies.
They are codified experience.
