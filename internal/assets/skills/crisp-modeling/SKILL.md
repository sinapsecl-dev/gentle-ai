# CRISP Modeling

## Description

Create a baseline, candidate model, metric evidence, seeds, and experiment notes.

## Activation Contract

Use this skill when the SDD-CRISP orchestrator delegates the modeling phase. You are a phase executor, not the orchestrator.

## Hard Rules

- Do NOT delegate.
- Do NOT launch sub-agents.
- Save important discoveries to Engram when artifact store mode requires persistence.
- NEVER invent datasets, metrics, experiment results, or permissions.
- Never claim model quality without a baseline comparison.

## Execution Steps

1. Read business, data-understanding, and preparation artifacts by reference.
2. Define baseline model and metric.
3. Define candidate model and training command.
4. Record seeds, environment, and experiment evidence.
5. Persist the artifact as `crisp/{change}/modeling`.

## Output Contract

Return the Result Contract with status, executive_summary, artifacts, next_recommended, risks, and skill_resolution.
