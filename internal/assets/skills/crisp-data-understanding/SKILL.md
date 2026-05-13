# CRISP Data Understanding

## Description

Map data sources, schemas, quality risks, permissions, PII, and freshness before preparation or modeling.

## Activation Contract

Use this skill when the SDD-CRISP orchestrator delegates the data-understanding phase. You are a phase executor, not the orchestrator.

## Hard Rules

- Do NOT delegate.
- Do NOT launch sub-agents.
- Use injected project standards before making recommendations.
- Save important discoveries to Engram when artifact store mode requires persistence.
- NEVER invent datasets, metrics, experiment results, or permissions.
- Treat unknown permissions or PII as risks, not assumptions.

## Execution Steps

1. Read `crisp/{change}/business-understanding` by reference.
2. Identify data sources and whether they are mock, local, Postgres, GCP, API, or unknown.
3. Record expected schema, target column, timestamp column, and key entities.
4. Assess quality risks: missingness, duplicates, outliers, freshness, and granularity.
5. Assess governance risks: PII, permissions, retention, and source authorization.
6. Persist the artifact as `crisp/{change}/data-understanding`.

## Output Contract

Return the Result Contract with status, executive_summary, artifacts, next_recommended, risks, and skill_resolution.
