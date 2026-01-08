# Contributing to production-readiness

Thank you for considering contributing.

This project values **operational judgment** over clever abstractions.

It grows through real operational experience.
If you have lived through incidents and want fewer of them, your input matters more than perfect code.

---

## What we welcome most

- New production-readiness rules
- New detectors for platforms and tools
- Real incident patterns
- Documentation improvements
- Test cases from real systems

---

## Contribution types

### 1. Rules

Add a new YAML file in `/rules`.

Checklist:

- Clear `id`
- Conservative detection
- Strong `why_it_matters`
- Vendor-neutral language

---

### 2. Detectors

Add detection logic in: `internal/scanner/`

Checklist:

- Deterministic
- No external calls
- No mutation
- Safe on large repos

---

## Philosophy for contributors

This project is not about:

- perfect coverage
- theoretical correctness
- enforcing compliance

It is about:

- surfacing **real-world risk**
- encoding **hard-earned lessons**
- helping teams fail less painfully

---

## How to propose a rule

When opening a PR, please include:

1. The rule YAML
2. A short incident story:
   - What went wrong?
   - Why existing tools didn’t catch it?
   - How this rule would have helped?

This keeps the rule set grounded in reality.

---

## Code style

- Keep logic simple
- Prefer clarity over abstraction
- Avoid premature generalization

This project optimizes for **maintainability**, not cleverness.

---

## Code of conduct

Be respectful.
Assume good intent.
Debate ideas, not people.
