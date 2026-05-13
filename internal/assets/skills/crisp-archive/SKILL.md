# CRISP Archive

## Description

Close the CRISP cycle by preserving final decisions, learnings, artifacts, model/data cards when applicable, and next actions.

## Activation Contract

Use this skill when the SDD-CRISP orchestrator delegates archive. You are a phase executor, not the orchestrator.

## Hard Rules

- Do NOT delegate.
- Do NOT launch sub-agents.
- Save archive reports to Engram when artifact store mode requires persistence.
- NEVER invent datasets, metrics, experiment results, or permissions.
- Do not archive as successful when verify has unresolved CRITICAL findings.

## Execution Steps

1. Read all CRISP artifacts and verify report by reference.
2. Summarize final decision, evidence, risks, and remaining work.
3. Record reusable learnings and artifact links.
4. Persist the artifact as `crisp/{change}/archive-report`.

## Output Contract

Return the Result Contract with status, executive_summary, artifacts, next_recommended, risks, and skill_resolution.
