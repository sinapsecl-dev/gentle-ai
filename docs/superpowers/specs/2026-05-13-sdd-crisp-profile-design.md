# SDD-CRISP Native Profile Design

Create a native SDD-CRISP profile that adapts Gentle-AI's existing SDD orchestration model to Data Science and ML product work. CRISP-DM becomes the visible domain workflow, while the current SDD architecture remains the internal discipline layer for artifacts, delegation, verification, memory, and review safety.

## Quick path

1. Keep the existing SDD flow intact for software/product engineering.
2. Add an SDD-CRISP profile with dedicated `crisp-*` phase agents that reuse the current SDD contracts.
3. Validate the MVP with a simple Python forecasting demo using mock data.
4. Verify structural and operational parity against the original SDD ecosystem before treating the profile as ready.

## Core decision

| Topic | Decision |
|-------|----------|
| Product shape | Add a native `sdd-crisp` profile, not a separate orchestration framework. |
| UX | Normal use is conversational. Commands may exist as expert/debug entrypoints. |
| Mental model | CRISP-DM is the Data Science facade; SDD is the internal chassis. |
| Compatibility | Existing `/sdd-*` behavior must remain intact. |
| First MVP | Simple forecasting flow with mock data. |
| Coexistence | Use a separate binary/config/profile strategy for the experimental fork. |

## Goals

- Guide Data Science and ML product work through CRISP-DM phases.
- Preserve the proven SDD structure: orchestrator, subagents, skills, artifact stores, Engram, result contracts, and verification gates.
- Support Python ML projects with strong defaults while adapting to existing repo conventions.
- Prevent hallucinated DS workflows by validating parity with the original SDD ecosystem.
- Keep the original installed Gentle-AI ecosystem easy to use while the fork is experimental.

## Non-goals

- Replace the existing SDD software-engineering flow.
- Build a full MLOps platform in the first release.
- Require real GitLab, Jenkins, GCP, Postgres, or API integrations for the MVP.
- Add complex `activate` / `deactivate` switching in the first version.
- Enforce heavy governance/compliance gates for normal exploration.

## Architecture

The profile extends the current SDD architecture instead of reinventing it.

```txt
gentle-orchestrator
  -> detects DS/ML intent
  -> activates SDD-CRISP profile internally
  -> delegates to crisp-* phase agents
  -> persists artifacts through the selected artifact store
  -> verifies reproducibility, metrics, governance, and readiness
```

### Hard rule

> SDD-CRISP must adapt existing SDD construction patterns. It must not introduce a parallel orchestration framework.

Implementation should treat these existing areas as source patterns:

- `internal/components/sdd/prompts.go`
- `internal/components/sdd/profiles.go`
- `internal/components/sdd/read_assignments.go`
- `.pi/agents/*`
- `.pi/chains/*`
- existing `skills/sdd-*` prompt and result contracts
- `.atl/skill-registry.md` compact-rule injection contract

## Agent model

### Orchestrator

The orchestrator remains a coordinator, not an executor.

Responsibilities:

- Detect whether the request is normal SDD or SDD-CRISP.
- Ensure SDD/crisp init artifacts exist.
- Read compact project standards once and inject matching rules into subagent prompts.
- Pass artifact references and minimal context, not large context dumps.
- Delegate phase work to `crisp-*` agents.
- Enforce workload/risk gates.
- Ask the user only at real decision points.

### Phase agents

Proposed agents:

```txt
crisp-business
crisp-data-understanding
crisp-data-preparation
crisp-modeling
crisp-evaluation
crisp-deployment
crisp-verify
crisp-archive
```

Every phase agent must:

- Execute one phase only.
- Not delegate or launch child agents.
- Read required artifacts by reference.
- Persist phase artifacts through the active backend.
- Save important discoveries/decisions to Engram when appropriate.
- Return the standard result envelope.
- Respect injected project standards.
- Detect and follow existing repo conventions before recommending new tools.

## Result contract

All SDD-CRISP phase agents return the same envelope used by SDD:

```txt
status
executive_summary
artifacts
next_recommended
risks
skill_resolution
```

## Artifact model

### Init artifact

```txt
crisp-init/{project}
```

Expected content:

```yaml
project_type: data-product | experiment | pipeline | api-ml-service
language: python
package_manager: uv | poetry | pip | conda | unknown
test_runner: pytest | unknown
lint_runner: ruff | unknown
type_checker: mypy | pyright | none | unknown
data_sources:
  - postgres
  - gcp
  - api
deployment:
  source_control: gitlab | unknown
  cd: jenkins | unknown
tracking:
  experiment: mlflow | custom | none | unknown
validation:
  data: pandera | great-expectations | custom | none | unknown
governance_mode: advisory | gated | strict
model_risk_level: low | medium | high | regulated
```

Defaults:

```yaml
governance_mode: advisory
model_risk_level: medium
```

### Phase artifacts

```txt
crisp/{change}/business-understanding
crisp/{change}/data-understanding
crisp/{change}/data-preparation
crisp/{change}/modeling
crisp/{change}/evaluation
crisp/{change}/deployment
crisp/{change}/verify-report
crisp/{change}/archive-report
```

## CRISP phases

| Phase | Output |
|-------|--------|
| Business understanding | Objective, hypothesis, stakeholder, success metric, risk level. |
| Data understanding | Sources, schema, quality, permissions, PII/freshness risks. |
| Data preparation | Transformations, split strategy, lineage, leakage checks. |
| Modeling | Baseline, candidate model, metrics, seeds, experiment evidence. |
| Evaluation | Baseline comparison, error analysis, threshold rationale, decision. |
| Deployment | Readiness, Jenkins/deploy contract when present, rollback, monitoring/drift plan. |
| Verify | Evidence audit across code, data, metrics, reproducibility, governance, and readiness. |
| Archive | Final decision, learnings, artifact links, model/data cards when applicable. |

## Python/ML stack strategy

Use opinionated defaults for greenfield work and adaptive detection for existing repos.

### Greenfield defaults

- Python
- `pytest`
- `ruff`
- `pandas` or `polars`
- `scikit-learn`
- `pydantic`
- `uv` or `poetry`
- optional `MLflow`
- optional `pandera` or Great Expectations
- Docker/Jenkins/GitLab awareness for future deployment work

### Brownfield detection

`crisp-init` should inspect common project files and adapt gates accordingly:

```txt
pyproject.toml
requirements.txt
poetry.lock
uv.lock
environment.yml
Makefile
tox.ini
noxfile.py
pytest.ini
mypy.ini
ruff.toml
Jenkinsfile
.gitlab-ci.yml
Dockerfile
notebooks/
src/
tests/
```

Rule:

> Detect, respect, then recommend. Do not replace existing stack choices without reason.

## Governance and model risk

Governance is explicit but progressive. It guides normal CRISP-DM flow without blocking exploration by default.

| Mode | Behavior |
|------|----------|
| `advisory` | Warns and records minimum evidence. Does not block normal progress. |
| `gated` | Blocks critical decisions such as deploy, PII use, or automated decisions. |
| `strict` | Requires formal evidence/approval for high-risk or regulated work. |

Cross-phase checks:

| Phase | Governance check |
|-------|------------------|
| Business | Impact, stakeholder, human decision boundary. |
| Data understanding | PII, permissions, authorized source, retention. |
| Preparation | Leakage, lineage, reproducibility. |
| Modeling | Bias, explainability need, experiment tracking. |
| Evaluation | Segment error analysis, threshold rationale. |
| Deployment | Monitoring, drift, rollback, approval. |

Principle:

> Do not block learning. Block irresponsible deployment.

## Verification model

`crisp-verify` asks a stronger question than “does it run?”

> Can we trust this result enough to move to the next phase?

Verification layers:

- Code quality: tests/lint/type checks based on detected stack.
- Data quality: schema, freshness, missingness, data validation evidence when available.
- Reproducibility: seeds, environment, commands, deterministic artifacts where practical.
- Experiment quality: baseline, metric, comparison, tracking evidence.
- Evaluation: error analysis, segment checks, acceptance decision.
- Governance: leakage, PII, bias/explainability according to risk level.
- Deployment readiness: Jenkins/deploy contract, rollback, monitoring/drift plan when applicable.

Severity levels:

| Severity | Example |
|----------|---------|
| CRITICAL | High leakage risk unresolved, no baseline, deploy without rollback. |
| WARNING | Missing segment analysis, incomplete tracking evidence. |
| SUGGESTION | Improve documentation, add model card, improve observability. |

## Coexistence strategy

For the first version, avoid complex activation switching. Use separate binary/config/profile isolation.

```txt
gentle-ai        -> original installed ecosystem
gentle-ai-crisp  -> experimental SDD-CRISP fork
```

The fork should use separate names where possible:

```txt
~/.gentle-ai-crisp
~/.config/gentle-ai-crisp
opencode profile: sdd-crisp
agent/profile names: crisp / sdd-crisp
```

Safety rule:

> The experimental fork must not silently modify the original installation/configuration. If shared config must be touched, the operation must be explicit and backed up.

Future option:

- Add `activate` / `deactivate` only after the separated flavor proves useful and safe.

## MVP demo

Validate the system with a simple Python forecasting project using mock data.

Example dataset:

```csv
date,sales,promo_flag,weekday
2026-01-01,100,0,4
2026-01-02,115,1,5
```

Demo goal:

```txt
Predict simple future sales from mock daily sales data.
```

Suggested stack:

- Python
- `pytest`
- `ruff`
- `pandas`
- `scikit-learn`
- local CSV mock data

The MVP should prove:

- `crisp-init` detects the demo stack.
- Business/data/preparation/modeling/evaluation artifacts are produced.
- A baseline exists.
- A simple model is trained/evaluated.
- MAE or MAPE is captured.
- Leakage checks are recorded.
- `crisp-verify` returns evidence-based PASS/WARNING/CRITICAL results.
- `crisp-archive` captures final learnings.

## Structural parity validation

Every SDD-CRISP component should be compared against the original SDD equivalent.

| Original SDD | SDD-CRISP | Validation |
|--------------|-----------|------------|
| Orchestrator rules | CRISP routing rules | Coordinates, does not execute. |
| `sdd-*` phase prompts | `crisp-*` phase prompts | Executes one phase, does not delegate. |
| Result contract | Result contract | Same required fields. |
| Artifact policy | CRISP artifact policy | Same backend modes. |
| Skill registry injection | CRISP injection | Compact rules present. |
| Apply/verify/archive flow | Modeling/verify/archive flow | Evidence and closure are traceable. |
| Workload guard | Workload/risk guard | Large/risky work requires decision. |

Parity rule:

> If SDD-CRISP deviates from SDD, the design must explain why. Otherwise, adapt the existing pattern.

## Operational parity matrix

These checks prevent “prompt similarity” from hiding behavioral drift.

| Area | Acceptance check |
|------|------------------|
| Tool-use contract | Orchestrator and subagents declare allowed/forbidden tool patterns. |
| MCP usage | Engram, Context7, Git/GitLab, and secret-handling rules are referenced where relevant. |
| Engram protocol | Init/phase/verify/archive artifacts use correct topic keys and persistence behavior. |
| Artifact stores | `engram`, `openspec`, `hybrid`, and `none` behavior mirrors SDD. |
| Skill usage | Relevant skills/compact rules are injected before phase execution. |
| Context delegation | Orchestrator passes references and compact context, not unnecessary broad dumps. |
| Subagent discipline | `crisp-*` agents do not delegate or orchestrate. |
| Never-do rules | Prompts forbid inventing datasets, metrics, results, or permissions. |
| Continuity | Progress artifacts are merged on continuation, not overwritten. |
| Verify evidence | PASS/WARNING/CRITICAL findings cite real artifacts/commands/evidence. |

### Never-do rules for SDD-CRISP

```txt
NEVER:
- invent datasets, metrics, experiment results, or business conclusions
- assume a model is good without a baseline
- use real PII or restricted data without confirmed permission
- hide leakage, data quality, or reproducibility risks
- replace an existing repo stack without detecting it first
- skip Engram/artifact persistence when the backend requires it
- pass to deployment without verification evidence
- let a subagent delegate or coordinate phases
- silently modify the original Gentle-AI installation/configuration
```

## Acceptance criteria

- The original SDD flow still works unchanged.
- `gentle-ai-crisp` can exist separately from `gentle-ai`.
- `crisp-init/{project}` captures Python/ML stack and governance defaults.
- All `crisp-*` phase agents return the standard result contract.
- Phase prompts include executor discipline, artifact policy, Engram behavior, and skill-resolution expectations.
- Structural parity tests compare `sdd-*` and `crisp-*` prompt requirements.
- Operational parity tests cover tools/MCP/Engram/context/never-do protocols.
- Forecasting demo completes with artifacts, baseline, model metric, verify report, and archive.
- Governance defaults to advisory/medium and does not block normal exploration.

## Open decisions for implementation planning

These are intentionally deferred until implementation planning:

- Exact binary/package rename strategy for `gentle-ai-crisp`.
- Whether `crisp-*` assets live under existing SDD component code or a sibling component with shared helpers.
- Exact golden-test fixture layout.
- Exact demo project location and generated fixture ownership.
- Whether MVP includes a local MLflow stub or leaves tracking as optional evidence.

## Next step

Review this design. If approved, create an implementation plan that slices work into small reviewable units and starts with prompt/profile parity before runtime routing or demo automation.
