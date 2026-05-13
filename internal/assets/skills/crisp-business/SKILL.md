# CRISP Business Understanding

## Description

Define the Data Science business objective, hypothesis, stakeholders, success metric, and risk level before data or modeling work begins.

## Activation Contract

Use this skill when the SDD-CRISP orchestrator delegates the business-understanding phase. You are a phase executor, not the orchestrator.

## Hard Rules

- Do NOT delegate.
- Do NOT launch sub-agents.
- Use injected project standards before making recommendations.
- Save important decisions to Engram when artifact store mode requires persistence.
- NEVER invent datasets, metrics, experiment results, or permissions.
- Do not move to data understanding without a business objective and success metric.

## Execution Steps

1. Read the CRISP init artifact when present: `crisp-init/{project}`.
2. Capture the business objective in one sentence.
3. Identify stakeholders and decision users.
4. Define the target prediction/analysis question.
5. Define one primary success metric and one fallback metric.
6. Set `governance_mode` and `model_risk_level` if missing.
7. Persist the artifact as `crisp/{change}/business-understanding`.

## Output Contract

Return the Result Contract with:

- `status`
- `executive_summary`
- `artifacts`
- `next_recommended`
- `risks`
- `skill_resolution`

The artifact must include objective, hypothesis, stakeholders, success metric, decision boundary, governance mode, and model risk level.
