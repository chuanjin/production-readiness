# Writing Detectors

Detectors are responsible for **finding signals** in a repository.  
Rules consume signals — they never inspect files directly.

---

## What is a signal?

A signal is a simple, deterministic fact about the system, such as:

- `secrets_provider_detected = true`
- `infra_as_code_detected = false`
- `correlation_id_detected = true`

Signals are stored in `RepoSignals`.

---

## Where detectors live

Under folder internal/scanner/

- `fs.go` → file system scanning
- `signals.go` → signal structure definition

---

## Supported signal maps

```go
 Files         map[string]bool   // tracks file existence
 FileContent   map[string]string // scanned file content (code only)
 BoolSignals   map[string]bool
 StringSignals map[string]string
 IntSignals    map[string]int

```

## Adding a new signal

1. Add the signal detect function from `detectors.go`

```go
func detectSecretsProvider(signals *RepoSignals) {
 if containsAny(signals.FileContent, "vault", "aws secretsmanager") {
  signals.BoolSignals["secrets_provider_detected"] = true
 }
}
```

2. Register the function in `init()`

```go
 registerDetector(detectSecretsProvider)

```

3. Use it in a rule

```yaml
detect:
  none_of:
    - signal_equals:
        secrets_provider_detected: true

```

## Detector design principles

- Deterministic
- No network calls
- No mutation
- No side effects
- Best-effort, not perfect
- Conservative in claiming presence

## What detectors should NOT do

- Deploy anything
- Modify files
- Enforce policy
- Make assumptions about intent

Their job is simple:
> “Is this signal likely present?”
