# CRISP Evaluation

## Description

Compare baseline and candidate results, analyze errors, explain acceptance decisions, and record threshold rationale.

## Activation Contract

Use this skill when the SDD-CRISP orchestrator delegates the evaluation phase. You are a phase executor, not the orchestrator.

## Hard Rules

- Do NOT delegate.
- Do NOT launch sub-agents.
- Save important decisions to Engram when artifact store mode requires persistence.
- NEVER invent datasets, metrics, experiment results, or permissions.
- Do not approve a model without metric evidence and known limitations.

## Execution Steps

1. Read modeling and preparation artifacts by reference.
2. Compare candidate model against baseline.
3. Record error analysis and relevant segments.
4. Decide whether to iterate, deploy, or stop.
5. Persist the artifact as `crisp/{change}/evaluation`.

## Output Contract

Return the Result Contract with status, executive_summary, artifacts, next_recommended, risks, and skill_resolution.
