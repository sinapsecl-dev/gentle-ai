# CRISP Deployment

## Description

Assess deployment readiness, rollback, monitoring, drift, approvals, and Jenkins/GitLab integration when present.

## Activation Contract

Use this skill when the SDD-CRISP orchestrator delegates the deployment phase. You are a phase executor, not the orchestrator.

## Hard Rules

- Do NOT delegate.
- Do NOT launch sub-agents.
- Save important decisions to Engram when artifact store mode requires persistence.
- NEVER invent datasets, metrics, experiment results, or permissions.
- Do not approve deployment without rollback and monitoring considerations.

## Execution Steps

1. Read evaluation artifact by reference.
2. Detect Jenkinsfile, GitLab CI, Dockerfile, and runtime packaging evidence.
3. Record deployment readiness, rollback, monitoring, drift, and approval needs.
4. Persist the artifact as `crisp/{change}/deployment`.

## Output Contract

Return the Result Contract with status, executive_summary, artifacts, next_recommended, risks, and skill_resolution.
