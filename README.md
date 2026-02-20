# production-readiness

**Turn senior engineering intuition into automated checks.**

> **This is not a scanner of code.
> This is a scanner of operational blind spots.**

Most systems don’t fail because of bugs.
They fail because they were never truly production-ready.

**production-readiness** is a read-only, opinionated tool that evaluates whether
a system is *actually safe to run in production* — based on the same mental
checklists senior engineers use when reviewing real systems before they
go live.

This project is for engineers who already run real systems in production and
want fewer surprises.

If you are responsible for availability, on-call, or launch decisions,
this tool is for you.

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

## What this catches in practice

In real systems, this tool typically surfaces issues like:

- A deployment pipeline that has no rollback path, even though rollbacks are assumed
- Database migrations that are not backward-compatible and will fail under load
- Services with metrics but no request correlation, making incidents hard to debug
- Rate limiting missing at the edge, leading to cascading failures
- Secrets drifting into environment files “temporarily” and never leaving
- Kubernetes workloads running without resource limits, risking node instability
- Lack of graceful shutdown handling, leading to dropped requests during deploys
- Missing SLO or Error Budget configurations for critical services
- Missing or inconsistent timeout and retry configurations

These are rarely flagged by linters or security scanners, but they are common causes of real production incidents.
If you have ever said “we should have seen this coming”, this tool is meant to make those risks visible earlier.

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

## Scope and non-goals

This project focuses on surfacing deterministic, explainable signals about production risk that are visible from code, configuration, and deployment intent.

It does not aim to:

- enumerate all possible runtime failure states of a system
- replace runtime testing, staging validation, or operational review
- predict incidents or guarantee correctness
- enforce best practices or auto-remediate changes

Many production failures only emerge under real traffic, timing, dependency behavior, or human interaction. Those require empirical validation and operational judgment.
This project is intentionally upstream of those activities and is meant to complement — not replace — existing engineering and operational practices.

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

1. **Scans**: Walk target repository files.
2. **Extracts**: Multiple specialized detectors extract production-readiness signals:
    - **Infrastructure Detector**: Scans IaC (Terraform, CloudFormation) and cloud provider patterns.
    - **Kubernetes Detector**: Evaluates Deployments, Ingress, and resource configurations.
    - **Reliability Detector**: Finds patterns for timeouts, retries, circuit breakers, and SLOs.
    - **Process Detector**: Looks for manual steps in documentation and migration patterns.
3. **Evaluates**: Correlates signals against a curated rule set.
4. **Reports**: Produces a summary of risks and maturity indicators.

For information about usage:

```
pr --help
```

Example output:

```
Overall Readiness Score: 62 / 100

🔴 High Risk
- No rollback strategy detected
- Secrets likely managed via environment variables
- Kubernetes workloads missing resource limits (CPU/Memory)

🟠 Medium Risk
- No rate limiting at ingress or API Gateway
- Logging without correlation IDs (Trace/Request ID)
- Missing Graceful Shutdown handling for SIGTERM

🟡 Low Risk
- No database migration safety signals (expand-contract)
- Service Level Objectives (SLO) not explicitly defined

🟢 Good Signals
- Health checks and Readiness probes detected
- Versioned deployment artifacts
- Infrastructure-as-Code (Terraform) detected

```

Each finding includes:

- what was detected
- why it matters in real incidents
- how teams usually get burned

### Rules

Rules live in rules/*.yaml and are fully open-source —
you can read, modify, or PR new ones.

Rules are intentionally opinionated,
reflecting common real-world failure patterns rather than theoretical best practices.

They are signals, not prescriptions.
A rule firing highlights an explicit assumption or trade-off,
not a required action or universal judgment.

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

## Extending production-readiness

`production-readiness` is designed to grow with real-world experience.

You can extend it by:

- Adding new production-readiness rules (YAML)
- Implementing new detectors for additional platforms and tools

Documentation:

- `docs/architecture.md` — system architecture and data flow
- `docs/rules.md` — how to write rules
- `docs/detectors.md` — how to add detectors
- `docs/contributing.md` — contribution guide

## Scope

This project focuses on **deterministic detection** of production-readiness signals.
Interpretation, workflow automation, and organizational policy are intentionally kept out of scope.

## Direction

### Short-term focus

- Expand detector coverage for Helm and more varied Terraform providers
- Improve report explanations with real incident patterns and "burn" stories
- Add language-specific detectors for more frameworks (Go, Node.js, Python)
- CI/CD integration guides (GitHub Actions, GitLab CI)

### Longer-term

- Keep this tool read-only and explainable
- Avoid turning it into a compliance or gatekeeping system
- Plugin architecture for custom detectors
- This project is meant to stay lightweight and opinionated.

## Star the project ⭐

If this reflects problems you have seen in production, a star helps signal that this direction is useful.

[![CI](https://github.com/chuanjin/production-readiness/actions/workflows/ci.yml/badge.svg)](https://github.com/chuanjin/production-readiness/actions/workflows/ci.yml)
[![codecov](https://codecov.io/gh/chuanjin/production-readiness/branch/main/graph/badge.svg)](https://codecov.io/gh/chuanjin/production-readiness)
[![Go Report Card](https://goreportcard.com/badge/github.com/chuanjin/production-readiness)](https://goreportcard.com/report/github.com/chuanjin/production-readiness)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
