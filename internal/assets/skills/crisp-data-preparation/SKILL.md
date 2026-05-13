# CRISP Data Preparation

## Description

Plan reproducible transformations, train/test split strategy, lineage, and leakage checks.

## Activation Contract

Use this skill when the SDD-CRISP orchestrator delegates the data-preparation phase. You are a phase executor, not the orchestrator.

## Hard Rules

- Do NOT delegate.
- Do NOT launch sub-agents.
- Save important discoveries to Engram when artifact store mode requires persistence.
- NEVER invent datasets, metrics, experiment results, or permissions.
- Never allow target leakage or future leakage to be hidden as a warning-only issue.

## Execution Steps

1. Read business and data-understanding artifacts by reference.
2. Define transformations and feature generation steps.
3. Define split strategy appropriate for the problem, using time-aware splits for forecasting.
4. Record leakage checks and lineage.
5. Persist the artifact as `crisp/{change}/data-preparation`.

## Output Contract

Return the Result Contract with status, executive_summary, artifacts, next_recommended, risks, and skill_resolution.
