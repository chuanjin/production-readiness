# production-readiness

Turn senior engineering intuition into automated checks.
A CLI that scans a repository and reports whether the system is truly ready for production — beyond CI passing, beyond tests being green.

Most failures in real systems are not bugs.
They are missing rollback paths, unsafe schema changes, secrets in env, no rate limits, no SLOs, single-region deployments, logs that cannot be debugged under pressure.

This tool encodes that experience into rules — and runs them automatically.

## Why this exists

Green CI does not mean safe to deploy.
Clean architecture does not mean resilient.
Monitoring does not mean operable.

Production-readiness is a hidden skill, learned only through outages.
This project makes it observable, testable, and teachable.

## What it does (v0.1 MVP)

* Scans code + config
* Loads rule definitions (YAML)
* Prints a readiness report (Markdown or JSON)
* No mutation, no deployment, no infrastructure access — read-only

## Install

From source

```
git clone <https://github.com/chuanjin/production-readiness>
cd production-readiness
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

## Rules

Rules live in rules/*.yaml and are fully open-source — you can read, modify, or PR new ones.
Rules are intentionally opinionated.
They reflect "what goes wrong in real world" rather than academia.

## Philosophy

This project is intentionally:

|Yes                            |No                             |
|-------------------------------|-----------------------------|
|Read-only (no deployment)       | A replacement for Terraform / Helm           |
| Opinionated                    | Vendor-neutral                               |
| Education-first                | Just another linter                          |
| Lightweight                    | Full enterprise platform                         |

Goal: provide signal, not control.
The output is meant for humans – CTOs, Tech Leads, Architects, SREs.

## Example Use Cases

* Tech Lead doing architecture review before approving deployment
* New joiner learning the system, teaches them “what matters”
* CTO reviewing vendors and compares readiness across repos

## Contributing

PRs welcome — especially:

* new rules (real-world failure stories welcome)
* rule packs for industries (FinTech, Healthcare, IoT)
* better scanners/detectors (Terraform, Helm, Kubernetes)

## Star the project ⭐

If this helps you — starring the repo helps visibility and keeps development going.
