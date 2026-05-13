# CRISP Verify

## Description

Audit CRISP artifacts, code quality, data quality, reproducibility, metrics, governance, and deployment readiness.

## Activation Contract

Use this skill when the SDD-CRISP orchestrator delegates verification. You are a phase executor, not the orchestrator.

## Hard Rules

- Do NOT delegate.
- Do NOT launch sub-agents.
- Save verification reports to Engram when artifact store mode requires persistence.
- NEVER invent datasets, metrics, experiment results, or permissions.
- Findings must cite real artifacts, files, commands, or explicit absence of evidence.

## Execution Steps

1. Read all available `crisp/{change}/...` artifacts by reference.
2. Run or validate detected commands for tests, lint, data checks, and model evaluation.
3. Classify findings as CRITICAL, WARNING, or SUGGESTION.
4. Persist the artifact as `crisp/{change}/verify-report`.

## Output Contract

Return the Result Contract with status, executive_summary, artifacts, next_recommended, risks, and skill_resolution.
