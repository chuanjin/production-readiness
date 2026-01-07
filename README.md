# production-readiness

**Turn senior engineering intuition into automated checks.**

> **This is not a scanner of code.
> This is a scanner of operational blind spots.**

Most systems don’t fail because of bugs.
They fail because they were never truly production-ready.

**production-readiness** is a read-only, opinionated tool that evaluates whether a system is *actually safe to run in production* — based on the same mental checklists senior engineers use when reviewing real systems before they go live.

---

## What this is

`production-readiness` scans:

- source code
- infrastructure-as-code
- CI/CD configuration
- deployment artifacts

...and produces a **Production Readiness Report** that highlights:

- high-risk operational gaps
- latent failure modes
- missing safety signals
- maturity indicators

It does **not** deploy anything.  
It does **not** enforce policy.  
It does **not** gate your pipeline (at least in current version).  

It only does one thing:

> **Tell you where your system is most likely to fail — and why.**

---

## Why this exists

Most teams already have:

- CI pipelines
- linters
- security scanners
- monitoring
- dashboards

And yet outages still happen.

Because incidents rarely come from what tools already check.
They come from what only experience sees:

- No rollback path
- Unsafe database migrations
- Missing rate limits
- One-region assumptions
- Secrets that are “temporarily” in env files
- Logging that looks fine until the incident

These are not *syntax* problems.
They are *operational design* problems.

This project exists to turn those invisible risks into visible signals.

---

## Why not just a checklist? Why not just AI?

Most companies already have production-readiness checklists.  
Most teams can ask AI for advice.  
Yet incidents keep happening.

Because:

- Checklists are **static** — systems are not.
- AI advice is **unbounded** — production risk is concrete.
- Human reviews are **inconsistent** — outages are not.

**production-readiness sits in the middle ground:**

| Checklists | AI | production-readiness |
|------------|----|---------------------|
| Static | Probabilistic | Deterministic |
| Manual | Unverifiable | Reproducible |
| Contextless | Context-heavy but vague | Context-aware and explicit |
| Forgotten after onboarding | Used only when asked | Run every time |

This tool turns **implicit expectations** into **executable standards**.

---

## Philosophy

This project is intentionally:

| Yes | No |
|-----|----|
| Opinionated about engineering outcomes | Opinionated about vendors |
| Read-only | Deployment or enforcement |
| Education-first | Compliance theater |
| Lightweight | Platform lock-in |

It behaves like a senior engineer reviewing a system before launch —
not like a tool enforcing policy after failure.

---

## How it works

### Install from source

```
git clone https://github.com/chuanjin/production-readiness
cd production-readiness
go mod tidy
go build -o pr ./cmd/pr
sudo mv pr /usr/local/bin
```

Run

```
pr scan .
```

or scan another repo:

```
pr scan ~/projects/my-microservice
```

The tool:

1. Scans the target repository
2. Extracts production-readiness signals
3. Evaluates them against a curated rule set
4. Produces a report in Markdown or JSON

For information about usage:

```
pr --help
```

Example output:

```
Overall Readiness Score: 68 / 100

🔴 High Risk
- No rollback strategy detected
- Secrets likely managed via environment variables

🟠 Medium Risk
- No rate limiting at ingress
- Logging without correlation IDs

🟡 Low Risk
- No database migration safety signals

🟢 Good Signals
- Health checks detected
- Versioned deployment artifacts

```

Each finding includes:

- what was detected
- why it matters in real incidents
- how teams usually get burned

### Rules

Rules live in rules/*.yaml and are fully open-source —
you can read, modify, or PR new ones.
Rules are intentionally opinionated.
They reflect "what goes wrong in real world" rather than academia.

---

## What this tool is NOT

This project is not:

- a CI/CD system
- a security scanner
- a Terraform validator
- a Kubernetes linter
- a compliance framework

It complements all of them by answering a different question:
> If this system fails in production, where will it most likely fail first?

---

## Who this is for

- Tech Leads
- Staff / Principal Engineers
- SREs / DevOps
- Startup founders shipping their first production system
- Teams that have already lived through outages and want fewer of them

Juniors use it to learn what seniors look for.
Seniors use it to scale their judgment.

### Example Use Cases

- Tech Lead doing architecture review before approving deployment
- New joiner learning the system, teaches them “what matters”
- CTO reviewing vendors and compares readiness across repos

## Contributing

PRs welcome — especially:

- new rules (real-world failure stories welcome)
- rule packs for industries (FinTech, MedTech, IoT)
- better scanners/detectors (Terraform, Helm, Kubernetes)

## Star the project ⭐

If this helps you — starring the repo helps visibility and keeps development going.

[![CI](https://github.com/chuanjin/production-readiness/actions/workflows/ci.yml/badge.svg)](https://github.com/chuanjin/production-readiness/actions/workflows/ci.yml)
[![codecov](https://codecov.io/gh/chuanjin/production-readiness/branch/main/graph/badge.svg)](https://codecov.io/gh/chuanjin/production-readiness)
[![Go Report Card](https://goreportcard.com/badge/github.com/chuanjin/production-readiness)](https://goreportcard.com/report/github.com/chuanjin/production-readiness)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
