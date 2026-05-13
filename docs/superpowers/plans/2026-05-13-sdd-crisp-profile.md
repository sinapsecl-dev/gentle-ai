# SDD-CRISP Native Profile Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Add a native SDD-CRISP profile that reuses the existing SDD orchestration architecture for Data Science/ML workflows and validates the MVP through a mock forecasting demo.

**Architecture:** Extend the existing `internal/components/sdd` patterns instead of creating a parallel framework. Add CRISP phase catalogs, prompt files, OpenCode profile overlay support, CRISP skill assets, parity tests, and a demo fixture while preserving the existing `/sdd-*` flow unchanged.

**Tech Stack:** Go 1.25, existing `internal/components/sdd` package, embedded assets, OpenCode agent overlay JSON, Engram artifact conventions, Python fixture files for the forecasting demo.

---

## Review path

Review in this order:

1. CRISP phase catalog and prompt writing.
2. OpenCode CRISP overlay/profile generation.
3. CRISP skill/orchestrator assets.
4. Structural and operational parity tests.
5. Mock forecasting demo fixture.
6. Separate `gentle-ai-crisp` binary/config strategy.

## File structure

### Create

- `internal/components/sdd/crisp.go` — CRISP phase constants, prompt directory, prompt content, and profile key helpers.
- `internal/components/sdd/crisp_test.go` — CRISP phase/prompt/profile-key tests.
- `internal/components/sdd/crisp_profile.go` — OpenCode CRISP profile overlay generation.
- `internal/components/sdd/crisp_profile_test.go` — overlay structure, permission, model assignment, and file-ref tests.
- `internal/assets/opencode/crisp-orchestrator.md` — OpenCode orchestrator addendum for CRISP routing.
- `internal/assets/skills/crisp-business/SKILL.md` — business understanding phase skill.
- `internal/assets/skills/crisp-data-understanding/SKILL.md` — data understanding phase skill.
- `internal/assets/skills/crisp-data-preparation/SKILL.md` — data preparation phase skill.
- `internal/assets/skills/crisp-modeling/SKILL.md` — modeling phase skill.
- `internal/assets/skills/crisp-evaluation/SKILL.md` — evaluation phase skill.
- `internal/assets/skills/crisp-deployment/SKILL.md` — deployment phase skill.
- `internal/assets/skills/crisp-verify/SKILL.md` — verification phase skill.
- `internal/assets/skills/crisp-archive/SKILL.md` — archive phase skill.
- `internal/components/sdd/crisp_assets_test.go` — validates CRISP skill/orchestrator asset markers.
- `internal/components/sdd/crisp_parity_test.go` — validates structural and operational parity against SDD patterns.
- `testdata/crisp-forecasting-demo/pyproject.toml` — Python demo project metadata.
- `testdata/crisp-forecasting-demo/data/daily_sales.csv` — mock sales data.
- `testdata/crisp-forecasting-demo/src/forecasting_demo/__init__.py` — demo Python package marker.
- `testdata/crisp-forecasting-demo/src/forecasting_demo/data.py` — load and validate mock data.
- `testdata/crisp-forecasting-demo/src/forecasting_demo/model.py` — baseline and simple forecasting model.
- `testdata/crisp-forecasting-demo/tests/test_forecasting.py` — demo tests.
- `internal/components/sdd/crisp_demo_test.go` — validates demo fixture files and expected commands.
- `cmd/gentle-ai-crisp/main.go` — experimental fork binary entrypoint.
- `internal/app/flavor.go` — runtime flavor metadata for binary name and config namespace.
- `internal/app/flavor_test.go` — flavor behavior tests.

### Modify

- `internal/components/sdd/prompts.go` — keep existing SDD functions unchanged; optionally call shared helpers if introduced during implementation.
- `internal/components/sdd/profiles.go` — keep existing SDD functions unchanged; reuse helper patterns for CRISP overlay without changing SDD behavior.
- `internal/components/sdd/inject.go` — wire CRISP prompt/profile generation behind explicit options only.
- `internal/components/skills/presets.go` — no MVP changes.
- `internal/model/types.go` — no MVP changes.
- `internal/app/app.go` — route binary name/help/version through flavor metadata without changing current `gentle-ai` behavior.
- `docs/superpowers/specs/2026-05-13-sdd-crisp-profile-design.md` — update only if implementation discovers a necessary design correction.

---

## Task 1: Add CRISP phase catalog and shared prompt files

**Files:**
- Create: `internal/components/sdd/crisp.go`
- Create: `internal/components/sdd/crisp_test.go`
- Do not modify: existing SDD phase order or `WriteSharedPromptFiles` behavior.

- [ ] **Step 1: Write the failing CRISP phase-order test**

Create `internal/components/sdd/crisp_test.go` with:

```go
package sdd

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestCRISPPhaseOrder(t *testing.T) {
	want := []string{
		"crisp-business",
		"crisp-data-understanding",
		"crisp-data-preparation",
		"crisp-modeling",
		"crisp-evaluation",
		"crisp-deployment",
		"crisp-verify",
		"crisp-archive",
	}

	got := CRISPPhaseOrder()
	if len(got) != len(want) {
		t.Fatalf("CRISPPhaseOrder() length = %d, want %d; got %v", len(got), len(want), got)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("CRISPPhaseOrder()[%d] = %q, want %q; got %v", i, got[i], want[i], got)
		}
	}

	got[0] = "mutated"
	if CRISPPhaseOrder()[0] != "crisp-business" {
		t.Fatal("CRISPPhaseOrder() returned mutable backing slice")
	}
}

func TestCRISPSharedPromptDir(t *testing.T) {
	want := filepath.FromSlash("/home/testuser/.config/opencode/prompts/crisp")
	got := CRISPSharedPromptDir(filepath.FromSlash("/home/testuser"))
	if got != want {
		t.Fatalf("CRISPSharedPromptDir() = %q, want %q", got, want)
	}
}

func TestWriteCRISPSharedPromptFilesCreatesFiles(t *testing.T) {
	home := t.TempDir()

	changed, err := WriteCRISPSharedPromptFiles(home)
	if err != nil {
		t.Fatalf("WriteCRISPSharedPromptFiles() error = %v", err)
	}
	if !changed {
		t.Fatal("WriteCRISPSharedPromptFiles() first call changed = false, want true")
	}

	for _, phase := range CRISPPhaseOrder() {
		path := filepath.Join(CRISPSharedPromptDir(home), phase+".md")
		data, err := os.ReadFile(path)
		if err != nil {
			t.Fatalf("expected CRISP prompt file %q: %v", path, err)
		}
		content := string(data)
		for _, marker := range []string{
			"not the orchestrator",
			"Do NOT delegate",
			"Do NOT call task/delegate",
			"Do NOT launch sub-agents",
			"Return the standard Result Contract",
			"Save important discoveries to Engram",
		} {
			if !strings.Contains(content, marker) {
				t.Fatalf("%s missing marker %q", phase, marker)
			}
		}
	}
}

func TestWriteCRISPSharedPromptFilesIdempotent(t *testing.T) {
	home := t.TempDir()
	first, err := WriteCRISPSharedPromptFiles(home)
	if err != nil {
		t.Fatalf("first WriteCRISPSharedPromptFiles() error = %v", err)
	}
	if !first {
		t.Fatal("first WriteCRISPSharedPromptFiles() changed = false")
	}
	second, err := WriteCRISPSharedPromptFiles(home)
	if err != nil {
		t.Fatalf("second WriteCRISPSharedPromptFiles() error = %v", err)
	}
	if second {
		t.Fatal("second WriteCRISPSharedPromptFiles() changed = true, want false")
	}
}
```

- [ ] **Step 2: Run the narrow test and verify it fails**

Run:

```powershell
go test ./internal/components/sdd -run TestCRISP -v
```

Expected: FAIL with undefined symbols such as `CRISPPhaseOrder`, `CRISPSharedPromptDir`, and `WriteCRISPSharedPromptFiles`.

- [ ] **Step 3: Implement the minimal CRISP prompt catalog**

Create `internal/components/sdd/crisp.go`:

```go
package sdd

import (
	"path/filepath"
	"strings"

	"github.com/gentleman-programming/gentle-ai/internal/components/filemerge"
)

var crispPhaseOrder = []string{
	"crisp-business",
	"crisp-data-understanding",
	"crisp-data-preparation",
	"crisp-modeling",
	"crisp-evaluation",
	"crisp-deployment",
	"crisp-verify",
	"crisp-archive",
}

var crispPhaseDescriptions = map[string]string{
	"crisp-business":          "Define business objective, hypothesis, stakeholder, success metric, and risk level",
	"crisp-data-understanding": "Map data sources, schema, quality, permissions, PII, and freshness risks",
	"crisp-data-preparation":   "Plan reproducible transformations, split strategy, lineage, and leakage checks",
	"crisp-modeling":           "Create baseline, candidate model, metric evidence, seeds, and experiment notes",
	"crisp-evaluation":         "Compare baseline/model results, error analysis, threshold rationale, and decision",
	"crisp-deployment":         "Assess deployment readiness, rollback, monitoring, drift, and approval needs",
	"crisp-verify":             "Audit CRISP artifacts, reproducibility, metrics, governance, and readiness",
	"crisp-archive":            "Close the CRISP cycle with decisions, learnings, artifact links, and cards when applicable",
}

func CRISPPhaseOrder() []string {
	return append([]string(nil), crispPhaseOrder...)
}

func CRISPSharedPromptDir(homeDir string) string {
	return filepath.Join(homeDir, ".config", "opencode", "prompts", "crisp")
}

func CRISPProfileAgentKeys(name string) []string {
	suffix := ""
	if name != "" {
		suffix = "-" + name
	}
	keys := make([]string, 0, len(crispPhaseOrder)+1)
	keys = append(keys, "crisp-orchestrator"+suffix)
	for _, phase := range crispPhaseOrder {
		keys = append(keys, phase+suffix)
	}
	return keys
}

func WriteCRISPSharedPromptFiles(homeDir string) (bool, error) {
	promptDir := CRISPSharedPromptDir(homeDir)
	anyChanged := false
	for _, phase := range crispPhaseOrder {
		content := crispSubAgentPrompt(phase)
		path := filepath.Join(promptDir, phase+".md")
		result, err := filemerge.WriteFileAtomic(path, []byte(content), 0o644)
		if err != nil {
			return false, err
		}
		if result.Changed {
			anyChanged = true
		}
	}
	return anyChanged, nil
}

func crispSubAgentPrompt(phase string) string {
	description := crispPhaseDescriptions[phase]
	skillName := phase
	return strings.Join([]string{
		"You are an SDD-CRISP executor for the " + phase + " phase, not the orchestrator.",
		"Phase purpose: " + description + ".",
		"Do this phase's work yourself.",
		"Do NOT delegate, Do NOT call task/delegate, and Do NOT launch sub-agents.",
		"Read your skill file at ~/.config/opencode/skills/" + skillName + "/SKILL.md and follow it exactly.",
		"Save important discoveries to Engram when the active artifact store or phase contract requires persistence.",
		"Return the standard Result Contract: status, executive_summary, artifacts, next_recommended, risks, skill_resolution.",
		"Never invent datasets, metrics, experiment results, permissions, or business conclusions.",
	}, " ")
}
```

- [ ] **Step 4: Run the CRISP catalog tests and verify they pass**

Run:

```powershell
go test ./internal/components/sdd -run TestCRISP -v
```

Expected: PASS.

- [ ] **Step 5: Run the existing SDD prompt tests to confirm no regression**

Run:

```powershell
go test ./internal/components/sdd -run TestWriteSharedPromptFiles -v
```

Expected: PASS.

- [ ] **Step 6: Commit this task**

```powershell
git add internal/components/sdd/crisp.go internal/components/sdd/crisp_test.go
git commit -m "feat: add crisp phase prompt catalog"
```

---

## Task 2: Add CRISP OpenCode profile overlay generation

**Files:**
- Create: `internal/components/sdd/crisp_profile.go`
- Create: `internal/components/sdd/crisp_profile_test.go`

- [ ] **Step 1: Write the failing overlay structure tests**

Create `internal/components/sdd/crisp_profile_test.go`:

```go
package sdd

import (
	"encoding/json"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gentleman-programming/gentle-ai/internal/model"
)

func makeCRISPProfile() model.Profile {
	haiku := model.ModelAssignment{ProviderID: "anthropic", ModelID: "claude-haiku-3-5"}
	assignments := map[string]model.ModelAssignment{}
	for _, phase := range CRISPPhaseOrder() {
		assignments[phase] = haiku
	}
	return model.Profile{
		Name:              "sdd-crisp",
		OrchestratorModel: haiku,
		PhaseAssignments:  assignments,
	}
}

func TestGenerateCRISPProfileOverlayStructure(t *testing.T) {
	home := t.TempDir()
	overlay, err := GenerateCRISPProfileOverlay(makeCRISPProfile(), home)
	if err != nil {
		t.Fatalf("GenerateCRISPProfileOverlay() error = %v", err)
	}

	var root map[string]any
	if err := json.Unmarshal(overlay, &root); err != nil {
		t.Fatalf("overlay is not valid JSON: %v", err)
	}
	agentMap := root["agent"].(map[string]any)
	if len(agentMap) != len(CRISPPhaseOrder())+1 {
		t.Fatalf("agent count = %d, want %d", len(agentMap), len(CRISPPhaseOrder())+1)
	}

	orch := agentMap["crisp-orchestrator-sdd-crisp"].(map[string]any)
	if mode := orch["mode"].(string); mode != "primary" {
		t.Fatalf("orchestrator mode = %q, want primary", mode)
	}
	if modelName := orch["model"].(string); modelName != "anthropic/claude-haiku-3-5" {
		t.Fatalf("orchestrator model = %q", modelName)
	}
	prompt := orch["prompt"].(string)
	for _, marker := range []string{"SDD-CRISP", "CRISP-DM", "coordinator", "not an executor", "crisp-business-sdd-crisp"} {
		if !strings.Contains(prompt, marker) {
			t.Fatalf("orchestrator prompt missing marker %q", marker)
		}
	}
}

func TestGenerateCRISPProfileOverlaySubagents(t *testing.T) {
	home := t.TempDir()
	overlay, err := GenerateCRISPProfileOverlay(makeCRISPProfile(), home)
	if err != nil {
		t.Fatalf("GenerateCRISPProfileOverlay() error = %v", err)
	}
	var root map[string]any
	if err := json.Unmarshal(overlay, &root); err != nil {
		t.Fatalf("overlay is not valid JSON: %v", err)
	}
	agentMap := root["agent"].(map[string]any)
	promptDir := CRISPSharedPromptDir(home)
	for _, phase := range CRISPPhaseOrder() {
		key := phase + "-sdd-crisp"
		agent := agentMap[key].(map[string]any)
		if mode := agent["mode"].(string); mode != "subagent" {
			t.Fatalf("%s mode = %q, want subagent", key, mode)
		}
		if hidden := agent["hidden"].(bool); !hidden {
			t.Fatalf("%s hidden = false, want true", key)
		}
		wantPrompt := "{file:" + filepath.Join(promptDir, phase+".md") + "}"
		if got := agent["prompt"].(string); got != wantPrompt {
			t.Fatalf("%s prompt = %q, want %q", key, got, wantPrompt)
		}
	}
}

func TestGenerateCRISPProfileOverlayTaskPermissions(t *testing.T) {
	home := t.TempDir()
	overlay, err := GenerateCRISPProfileOverlay(makeCRISPProfile(), home)
	if err != nil {
		t.Fatalf("GenerateCRISPProfileOverlay() error = %v", err)
	}
	var root map[string]any
	if err := json.Unmarshal(overlay, &root); err != nil {
		t.Fatalf("overlay is not valid JSON: %v", err)
	}
	orch := root["agent"].(map[string]any)["crisp-orchestrator-sdd-crisp"].(map[string]any)
	taskWrapper := orch["permission"].(map[string]any)["task"].(map[string]any)
	taskMap := taskWrapper["__replace__"].(map[string]any)
	if taskMap["*"] != "deny" {
		t.Fatalf("permission wildcard = %v, want deny", taskMap["*"])
	}
	for _, phase := range CRISPPhaseOrder() {
		key := phase + "-sdd-crisp"
		if taskMap[key] != "allow" {
			t.Fatalf("permission for %s = %v, want allow", key, taskMap[key])
		}
	}
	if taskMap["sdd-apply"] != nil {
		t.Fatalf("CRISP orchestrator unexpectedly allowed sdd-apply")
	}
}
```

- [ ] **Step 2: Run the overlay tests and verify they fail**

Run:

```powershell
go test ./internal/components/sdd -run TestGenerateCRISPProfileOverlay -v
```

Expected: FAIL with `undefined: GenerateCRISPProfileOverlay`.

- [ ] **Step 3: Implement CRISP profile overlay generation**

Create `internal/components/sdd/crisp_profile.go`:

```go
package sdd

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/gentleman-programming/gentle-ai/internal/assets"
	"github.com/gentleman-programming/gentle-ai/internal/model"
)

func GenerateCRISPProfileOverlay(profile model.Profile, homeDir string) ([]byte, error) {
	if profile.Name == "" || profile.Name == "default" {
		return nil, fmt.Errorf("GenerateCRISPProfileOverlay: profile name must be non-empty and not 'default'")
	}
	suffix := "-" + profile.Name
	agentMap := make(map[string]any, len(crispPhaseOrder)+1)

	taskPerms := map[string]any{"*": "deny"}
	for _, phase := range crispPhaseOrder {
		taskPerms[phase+suffix] = "allow"
	}

	orchPrompt := buildCRISPOrchestratorPrompt(profile)
	orchEntry := map[string]any{
		"mode":        "primary",
		"description": "SDD-CRISP Orchestrator (" + profile.Name + " profile) - coordinates CRISP-DM sub-agents, never does phase work inline",
		"prompt":      orchPrompt,
		"permission": map[string]any{
			"task": map[string]any{"__replace__": taskPerms},
		},
		"tools": map[string]any{
			"read":            true,
			"write":           true,
			"edit":            true,
			"bash":            true,
			"delegate":        true,
			"delegation_read": true,
			"delegation_list": true,
		},
	}
	if profile.OrchestratorModel.ProviderID != "" && profile.OrchestratorModel.ModelID != "" {
		orchEntry["model"] = profile.OrchestratorModel.FullID()
	}
	agentMap["crisp-orchestrator"+suffix] = orchEntry

	promptDir := CRISPSharedPromptDir(homeDir)
	for _, phase := range crispPhaseOrder {
		entry := map[string]any{
			"mode":        "subagent",
			"hidden":      true,
			"description": crispPhaseDescriptions[phase],
			"prompt":      "{file:" + filepath.Join(promptDir, phase+".md") + "}",
			"tools": map[string]any{
				"read":  true,
				"write": true,
				"edit":  true,
				"bash":  true,
			},
		}
		if assignment, ok := profile.PhaseAssignments[phase]; ok && assignment.ProviderID != "" && assignment.ModelID != "" {
			entry["model"] = assignment.FullID()
		}
		agentMap[phase+suffix] = entry
	}

	overlay := map[string]any{"agent": agentMap}
	out, err := json.MarshalIndent(overlay, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("marshal CRISP profile overlay: %w", err)
	}
	return append(out, '\n'), nil
}

func buildCRISPOrchestratorPrompt(profile model.Profile) string {
	base := assets.MustRead("opencode/crisp-orchestrator.md")
	suffix := "-" + profile.Name
	for _, phase := range crispPhaseOrder {
		base = strings.ReplaceAll(base, phase, phase+suffix)
	}
	return base
}
```

- [ ] **Step 4: Add the minimal CRISP orchestrator asset required by the test**

Create `internal/assets/opencode/crisp-orchestrator.md`:

```markdown
# SDD-CRISP Orchestrator

You are the SDD-CRISP coordinator, not an executor. CRISP-DM is the visible Data Science workflow; SDD remains the internal discipline for artifacts, Engram persistence, skill injection, verification, and review safety.

## Routing

Use SDD-CRISP when the user asks for Data Science, ML, forecasting, model evaluation, datasets, feature engineering, experiment tracking, leakage analysis, or deployment readiness for data products. Use normal SDD for ordinary software changes.

## Delegation

Delegate phase work to:

- crisp-business
- crisp-data-understanding
- crisp-data-preparation
- crisp-modeling
- crisp-evaluation
- crisp-deployment
- crisp-verify
- crisp-archive

You coordinate. You do not execute phase work inline unless delegation is unavailable. Pass compact context and artifact references, not large dumps.

## Protocols

- Resolve and inject project standards before launching subagents.
- Use Engram topic keys: `crisp-init/{project}` and `crisp/{change}/{phase}`.
- Preserve the Result Contract: status, executive_summary, artifacts, next_recommended, risks, skill_resolution.
- Never invent datasets, metrics, experiment results, permissions, or business conclusions.
```

- [ ] **Step 5: Run overlay tests and verify they pass**

Run:

```powershell
go test ./internal/components/sdd -run TestGenerateCRISPProfileOverlay -v
```

Expected: PASS.

- [ ] **Step 6: Commit this task**

```powershell
git add internal/components/sdd/crisp_profile.go internal/components/sdd/crisp_profile_test.go internal/assets/opencode/crisp-orchestrator.md
git commit -m "feat: add crisp opencode profile overlay"
```

---

## Task 3: Add CRISP skill assets with operational protocol parity

**Files:**
- Create: `internal/assets/skills/crisp-business/SKILL.md`
- Create: `internal/assets/skills/crisp-data-understanding/SKILL.md`
- Create: `internal/assets/skills/crisp-data-preparation/SKILL.md`
- Create: `internal/assets/skills/crisp-modeling/SKILL.md`
- Create: `internal/assets/skills/crisp-evaluation/SKILL.md`
- Create: `internal/assets/skills/crisp-deployment/SKILL.md`
- Create: `internal/assets/skills/crisp-verify/SKILL.md`
- Create: `internal/assets/skills/crisp-archive/SKILL.md`
- Create: `internal/components/sdd/crisp_assets_test.go`

- [ ] **Step 1: Write the failing asset marker test**

Create `internal/components/sdd/crisp_assets_test.go`:

```go
package sdd

import (
	"strings"
	"testing"

	"github.com/gentleman-programming/gentle-ai/internal/assets"
)

func TestCRISPSkillAssetsContainProtocolMarkers(t *testing.T) {
	markers := []string{
		"## Activation Contract",
		"## Hard Rules",
		"## Execution Steps",
		"## Output Contract",
		"Result Contract",
		"Engram",
		"Do NOT delegate",
		"NEVER invent datasets, metrics, experiment results, or permissions",
	}
	for _, phase := range CRISPPhaseOrder() {
		path := "skills/" + phase + "/SKILL.md"
		content := assets.MustRead(path)
		for _, marker := range markers {
			if !strings.Contains(content, marker) {
				t.Fatalf("%s missing marker %q", path, marker)
			}
		}
	}
}
```

- [ ] **Step 2: Run the asset test and verify it fails**

Run:

```powershell
go test ./internal/components/sdd -run TestCRISPSkillAssetsContainProtocolMarkers -v
```

Expected: FAIL because the CRISP skill assets do not exist.

- [ ] **Step 3: Create `crisp-business` skill asset**

Create `internal/assets/skills/crisp-business/SKILL.md`:

```markdown
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
```

- [ ] **Step 4: Create the remaining phase assets with phase-specific execution steps**

Create the seven remaining files with the same sections and these phase-specific required outputs:

`internal/assets/skills/crisp-data-understanding/SKILL.md`:

```markdown
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
```

`internal/assets/skills/crisp-data-preparation/SKILL.md`:

```markdown
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
```

`internal/assets/skills/crisp-modeling/SKILL.md`:

```markdown
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
```

`internal/assets/skills/crisp-evaluation/SKILL.md`:

```markdown
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
```

`internal/assets/skills/crisp-deployment/SKILL.md`:

```markdown
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
```

`internal/assets/skills/crisp-verify/SKILL.md`:

```markdown
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
```

`internal/assets/skills/crisp-archive/SKILL.md`:

```markdown
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
```

- [ ] **Step 5: Run the CRISP asset test**

Run:

```powershell
go test ./internal/components/sdd -run TestCRISPSkillAssetsContainProtocolMarkers -v
```

Expected: PASS.

- [ ] **Step 6: Commit this task**

```powershell
git add internal/assets/skills/crisp-* internal/components/sdd/crisp_assets_test.go
git commit -m "feat: add crisp phase skill assets"
```

---

## Task 4: Add structural and operational parity tests

**Files:**
- Create: `internal/components/sdd/crisp_parity_test.go`

- [ ] **Step 1: Write parity tests**

Create `internal/components/sdd/crisp_parity_test.go`:

```go
package sdd

import (
	"strings"
	"testing"
)

func TestCRISPPromptsMaintainExecutorDiscipline(t *testing.T) {
	for _, phase := range CRISPPhaseOrder() {
		content := crispSubAgentPrompt(phase)
		for _, marker := range []string{
			"not the orchestrator",
			"Do NOT delegate",
			"Do NOT call task/delegate",
			"Do NOT launch sub-agents",
			"Result Contract",
			"Engram",
		} {
			if !strings.Contains(content, marker) {
				t.Fatalf("%s prompt missing operational marker %q", phase, marker)
			}
		}
	}
}

func TestCRISPProfileKeysDoNotCollideWithSDDKeys(t *testing.T) {
	sddKeys := map[string]bool{}
	for _, key := range ProfileAgentKeys("sdd-crisp") {
		sddKeys[key] = true
	}
	for _, key := range CRISPProfileAgentKeys("sdd-crisp") {
		if sddKeys[key] {
			t.Fatalf("CRISP key %q collides with SDD key", key)
		}
	}
}

func TestCRISPPhaseCountDocumentsIntentionalDifference(t *testing.T) {
	if got := len(CRISPPhaseOrder()); got != 8 {
		t.Fatalf("CRISP phase count = %d, want 8", got)
	}
	if got := len(ProfilePhaseOrder()); got != 10 {
		t.Fatalf("SDD phase count = %d, want 10", got)
	}
}
```

- [ ] **Step 2: Run parity tests**

Run:

```powershell
go test ./internal/components/sdd -run TestCRISP.*Parity|TestCRISPPrompts|TestCRISPProfileKeys|TestCRISPPhaseCount -v
```

PowerShell note: if regex quoting is problematic, run:

```powershell
go test ./internal/components/sdd -run TestCRISP -v
```

Expected: PASS after Tasks 1-3 are complete.

- [ ] **Step 3: Commit this task**

```powershell
git add internal/components/sdd/crisp_parity_test.go
git commit -m "test: add crisp operational parity checks"
```

---

## Task 5: Wire CRISP assets behind explicit injection options

**Files:**
- Modify: `internal/components/sdd/inject.go`
- Modify: `internal/components/sdd/inject_test.go`

- [ ] **Step 1: Write failing injection option test**

Append to `internal/components/sdd/inject_test.go`:

```go
func TestInjectOpenCodeWithCRISPProfileWritesCRISPPromptsAndOverlay(t *testing.T) {
	home := t.TempDir()
	mockNoPackageManager(t)

	profile := model.Profile{Name: "sdd-crisp", PhaseAssignments: map[string]model.ModelAssignment{}}
	result, err := Inject(home, opencodeAdapter(), model.SDDModeMulti, InjectOptions{
		CRISPProfiles: []model.Profile{profile},
	})
	if err != nil {
		t.Fatalf("Inject(opencode, CRISPProfiles) error = %v", err)
	}
	if !result.Changed {
		t.Fatal("Inject(opencode, CRISPProfiles) changed = false")
	}

	for _, phase := range CRISPPhaseOrder() {
		path := filepath.Join(CRISPSharedPromptDir(home), phase+".md")
		if _, err := os.Stat(path); err != nil {
			t.Fatalf("expected CRISP prompt %q: %v", path, err)
		}
	}

	settingsPath := filepath.Join(home, ".config", "opencode", "opencode.json")
	content, err := os.ReadFile(settingsPath)
	if err != nil {
		t.Fatalf("ReadFile(opencode.json) error = %v", err)
	}
	text := string(content)
	for _, want := range []string{"crisp-orchestrator-sdd-crisp", "crisp-business-sdd-crisp", "crisp-verify-sdd-crisp"} {
		if !strings.Contains(text, want) {
			t.Fatalf("opencode.json missing %q", want)
		}
	}
	if !strings.Contains(text, "sdd-orchestrator") {
		t.Fatal("existing SDD orchestrator missing after CRISP injection")
	}
}
```

- [ ] **Step 2: Run the new injection test and verify it fails**

Run:

```powershell
go test ./internal/components/sdd -run TestInjectOpenCodeWithCRISPProfileWritesCRISPPromptsAndOverlay -v
```

Expected: FAIL because `InjectOptions.CRISPProfiles` does not exist.

- [ ] **Step 3: Add explicit CRISP injection option**

Modify `InjectOptions` in `internal/components/sdd/inject.go`:

```go
type InjectOptions struct {
	OpenCodeModelAssignments map[string]model.ModelAssignment
	ClaudeModelAssignments   map[string]model.ClaudeModelAlias
	KiroModelAssignments     map[string]model.ClaudeModelAlias
	WorkspaceDir             string
	StrictTDD                bool
	Profiles                 []model.Profile
	CRISPProfiles            []model.Profile
	PreserveOpenCodeOrchestratorPrompt bool
}
```

Keep field order readable. If gofmt aligns the final field differently, accept gofmt output.

- [ ] **Step 4: Wire CRISP prompt writing and overlay merge for OpenCode only**

In `Inject`, after the existing OpenCode multi-mode prompt/overlay work has completed, add an explicit guarded block:

```go
if adapter.Agent() == model.AgentOpenCode && len(opts.CRISPProfiles) > 0 {
	changedPrompts, err := WriteCRISPSharedPromptFiles(homeDir)
	if err != nil {
		return InjectionResult{}, err
	}
	changed = changed || changedPrompts
	files = append(files, CRISPSharedPromptDir(homeDir))

	settingsPath := filepath.Join(adapter.GlobalConfigDir(homeDir), "opencode.json")
	for _, profile := range opts.CRISPProfiles {
		overlay, err := GenerateCRISPProfileOverlay(profile, homeDir)
		if err != nil {
			return InjectionResult{}, err
		}
    mergeResult, err := mergeJSONFile(settingsPath, overlay)
		if err != nil {
			return InjectionResult{}, err
		}
		changed = changed || mergeResult.writeResult.Changed
		files = append(files, settingsPath)
	}
}
```

- [ ] **Step 5: Run the CRISP injection test**

Run:

```powershell
go test ./internal/components/sdd -run TestInjectOpenCodeWithCRISPProfileWritesCRISPPromptsAndOverlay -v
```

Expected: PASS.

- [ ] **Step 6: Run existing OpenCode injection tests for regression**

Run:

```powershell
go test ./internal/components/sdd -run TestInjectOpenCode -v
```

Expected: PASS.

- [ ] **Step 7: Commit this task**

```powershell
git add internal/components/sdd/inject.go internal/components/sdd/inject_test.go
git commit -m "feat: wire crisp opencode profile injection"
```

---

## Task 6: Add Python forecasting demo fixture

**Files:**
- Create: `testdata/crisp-forecasting-demo/pyproject.toml`
- Create: `testdata/crisp-forecasting-demo/data/daily_sales.csv`
- Create: `testdata/crisp-forecasting-demo/src/forecasting_demo/__init__.py`
- Create: `testdata/crisp-forecasting-demo/src/forecasting_demo/data.py`
- Create: `testdata/crisp-forecasting-demo/src/forecasting_demo/model.py`
- Create: `testdata/crisp-forecasting-demo/tests/test_forecasting.py`
- Create: `internal/components/sdd/crisp_demo_test.go`

- [ ] **Step 1: Write fixture validation test**

Create `internal/components/sdd/crisp_demo_test.go`:

```go
package sdd

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestCRISPForecastingDemoFixtureExists(t *testing.T) {
	root := filepath.Join("..", "..", "..", "testdata", "crisp-forecasting-demo")
	paths := []string{
		"pyproject.toml",
		filepath.Join("data", "daily_sales.csv"),
		filepath.Join("src", "forecasting_demo", "__init__.py"),
		filepath.Join("src", "forecasting_demo", "data.py"),
		filepath.Join("src", "forecasting_demo", "model.py"),
		filepath.Join("tests", "test_forecasting.py"),
	}
	for _, rel := range paths {
		if _, err := os.Stat(filepath.Join(root, rel)); err != nil {
			t.Fatalf("missing demo fixture file %s: %v", rel, err)
		}
	}
}

func TestCRISPForecastingDemoDocumentsMetrics(t *testing.T) {
	root := filepath.Join("..", "..", "..", "testdata", "crisp-forecasting-demo")
	data, err := os.ReadFile(filepath.Join(root, "tests", "test_forecasting.py"))
	if err != nil {
		t.Fatalf("read test_forecasting.py: %v", err)
	}
	text := string(data)
	for _, want := range []string{"baseline", "MAE", "candidate", "assert"} {
		if !strings.Contains(text, want) {
			t.Fatalf("demo tests missing %q", want)
		}
	}
}
```

- [ ] **Step 2: Run fixture test and verify it fails**

Run:

```powershell
go test ./internal/components/sdd -run TestCRISPForecastingDemo -v
```

Expected: FAIL because fixture files do not exist.

- [ ] **Step 3: Create `pyproject.toml`**

Create `testdata/crisp-forecasting-demo/pyproject.toml`:

```toml
[project]
name = "crisp-forecasting-demo"
version = "0.1.0"
requires-python = ">=3.11"
dependencies = [
  "pandas>=2.2",
  "scikit-learn>=1.5",
]

[project.optional-dependencies]
dev = [
  "pytest>=8.0",
  "ruff>=0.6",
]

[tool.pytest.ini_options]
pythonpath = ["src"]

[tool.ruff]
line-length = 100
```

- [ ] **Step 4: Create mock data**

Create `testdata/crisp-forecasting-demo/data/daily_sales.csv`:

```csv
date,sales,promo_flag,weekday
2026-01-01,100,0,4
2026-01-02,115,1,5
2026-01-03,98,0,6
2026-01-04,105,0,0
2026-01-05,120,1,1
2026-01-06,125,1,2
2026-01-07,118,0,3
2026-01-08,130,1,4
2026-01-09,128,0,5
2026-01-10,122,0,6
```

- [ ] **Step 5: Create Python package files**

Create `testdata/crisp-forecasting-demo/src/forecasting_demo/__init__.py`:

```python
"""Mock forecasting demo used to validate SDD-CRISP flow."""
```

Create `testdata/crisp-forecasting-demo/src/forecasting_demo/data.py`:

```python
from pathlib import Path

import pandas as pd


REQUIRED_COLUMNS = {"date", "sales", "promo_flag", "weekday"}


def load_sales(path: str | Path) -> pd.DataFrame:
    frame = pd.read_csv(path, parse_dates=["date"])
    missing = REQUIRED_COLUMNS.difference(frame.columns)
    if missing:
        raise ValueError(f"missing required columns: {sorted(missing)}")
    if frame["sales"].isna().any():
        raise ValueError("sales contains missing values")
    return frame.sort_values("date").reset_index(drop=True)
```

Create `testdata/crisp-forecasting-demo/src/forecasting_demo/model.py`:

```python
from dataclasses import dataclass

import pandas as pd
from sklearn.dummy import DummyRegressor
from sklearn.ensemble import RandomForestRegressor
from sklearn.metrics import mean_absolute_error


FEATURES = ["promo_flag", "weekday"]


@dataclass(frozen=True)
class ForecastResult:
    baseline_mae: float
    candidate_mae: float


def time_split(frame: pd.DataFrame, holdout: int = 3) -> tuple[pd.DataFrame, pd.DataFrame]:
    if len(frame) <= holdout:
        raise ValueError("not enough rows for holdout split")
    return frame.iloc[:-holdout].copy(), frame.iloc[-holdout:].copy()


def evaluate_forecast(frame: pd.DataFrame) -> ForecastResult:
    train, test = time_split(frame)
    x_train = train[FEATURES]
    y_train = train["sales"]
    x_test = test[FEATURES]
    y_test = test["sales"]

    baseline = DummyRegressor(strategy="mean")
    baseline.fit(x_train, y_train)
    baseline_pred = baseline.predict(x_test)

    candidate = RandomForestRegressor(n_estimators=20, random_state=42)
    candidate.fit(x_train, y_train)
    candidate_pred = candidate.predict(x_test)

    return ForecastResult(
        baseline_mae=float(mean_absolute_error(y_test, baseline_pred)),
        candidate_mae=float(mean_absolute_error(y_test, candidate_pred)),
    )
```

- [ ] **Step 6: Create Python demo tests**

Create `testdata/crisp-forecasting-demo/tests/test_forecasting.py`:

```python
from pathlib import Path

from forecasting_demo.data import load_sales
from forecasting_demo.model import evaluate_forecast, time_split


DATA_PATH = Path(__file__).resolve().parents[1] / "data" / "daily_sales.csv"


def test_load_sales_validates_required_columns():
    frame = load_sales(DATA_PATH)
    assert list(frame.columns) == ["date", "sales", "promo_flag", "weekday"]
    assert frame["sales"].isna().sum() == 0


def test_time_split_keeps_future_rows_for_holdout():
    frame = load_sales(DATA_PATH)
    train, test = time_split(frame, holdout=3)
    assert train["date"].max() < test["date"].min()
    assert len(test) == 3


def test_candidate_model_reports_mae_against_baseline():
    frame = load_sales(DATA_PATH)
    result = evaluate_forecast(frame)
    assert result.baseline_mae >= 0
    assert result.candidate_mae >= 0
    assert result.candidate_mae < result.baseline_mae * 2
```

- [ ] **Step 7: Run Go fixture tests**

Run:

```powershell
go test ./internal/components/sdd -run TestCRISPForecastingDemo -v
```

Expected: PASS.

- [ ] **Step 8: Run Python demo tests when Python dependencies are available**

Run from `testdata/crisp-forecasting-demo`:

```powershell
python -m pytest -q
```

Expected: `3 passed`. If dependencies are not installed, record the missing dependency error and do not alter Go code to hide it.

- [ ] **Step 9: Commit this task**

```powershell
git add testdata/crisp-forecasting-demo internal/components/sdd/crisp_demo_test.go
git commit -m "test: add crisp forecasting demo fixture"
```

---

## Task 7: Add experimental `gentle-ai-crisp` binary flavor metadata

**Files:**
- Create: `cmd/gentle-ai-crisp/main.go`
- Create: `internal/app/flavor.go`
- Create: `internal/app/flavor_test.go`
- Modify: `internal/app/app.go`

- [ ] **Step 1: Write failing app flavor tests**

Create `internal/app/flavor_test.go`:

```go
package app

import (
	"bytes"
	"strings"
	"testing"
)

func TestRunArgsUsesConfiguredBinaryNameForVersion(t *testing.T) {
	original := RuntimeFlavor
	RuntimeFlavor = Flavor{Name: "gentle-ai-crisp", ConfigNamespace: "gentle-ai-crisp"}
	t.Cleanup(func() { RuntimeFlavor = original })

	var out bytes.Buffer
	if err := RunArgs([]string{"version"}, &out); err != nil {
		t.Fatalf("RunArgs(version) error = %v", err)
	}
	if !strings.Contains(out.String(), "gentle-ai-crisp") {
		t.Fatalf("version output = %q, want gentle-ai-crisp", out.String())
	}
}

func TestDefaultRuntimeFlavorIsGentleAI(t *testing.T) {
	if RuntimeFlavor.Name != "gentle-ai" {
		t.Fatalf("RuntimeFlavor.Name = %q, want gentle-ai", RuntimeFlavor.Name)
	}
	if RuntimeFlavor.ConfigNamespace != "gentle-ai" {
		t.Fatalf("RuntimeFlavor.ConfigNamespace = %q, want gentle-ai", RuntimeFlavor.ConfigNamespace)
	}
}
```

- [ ] **Step 2: Run flavor tests and verify they fail**

Run:

```powershell
go test ./internal/app -run Test.*RuntimeFlavor|TestRunArgsUsesConfiguredBinaryName -v
```

Expected: FAIL because `RuntimeFlavor` and `Flavor` do not exist.

- [ ] **Step 3: Add flavor metadata**

Create `internal/app/flavor.go`:

```go
package app

type Flavor struct {
	Name            string
	ConfigNamespace string
}

var RuntimeFlavor = Flavor{Name: "gentle-ai", ConfigNamespace: "gentle-ai"}

func SetRuntimeFlavor(flavor Flavor) {
	if flavor.Name == "" {
		flavor.Name = "gentle-ai"
	}
	if flavor.ConfigNamespace == "" {
		flavor.ConfigNamespace = flavor.Name
	}
	RuntimeFlavor = flavor
}
```

- [ ] **Step 4: Update version/help/error strings to use runtime flavor name**

In `internal/app/app.go`, change:

```go
_, _ = fmt.Fprintf(stdout, "gentle-ai %s\n", Version)
```

to:

```go
_, _ = fmt.Fprintf(stdout, "%s %s\n", RuntimeFlavor.Name, Version)
```

Change the unknown command error:

```go
return fmt.Errorf("unknown command %q — run 'gentle-ai help' for available commands", args[0])
```

to:

```go
return fmt.Errorf("unknown command %q — run '%s help' for available commands", args[0], RuntimeFlavor.Name)
```

If `printHelp` hardcodes `gentle-ai`, update its usage text to receive or read `RuntimeFlavor.Name`.

- [ ] **Step 5: Add experimental binary entrypoint**

Create `cmd/gentle-ai-crisp/main.go`:

```go
package main

import (
	"fmt"
	"os"

	"github.com/gentleman-programming/gentle-ai/internal/app"
)

var version = "dev"

func main() {
	app.Version = version
	app.SetRuntimeFlavor(app.Flavor{Name: "gentle-ai-crisp", ConfigNamespace: "gentle-ai-crisp"})
	if err := app.Run(); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
```

- [ ] **Step 6: Run app tests**

Run:

```powershell
go test ./internal/app -run Test.*RuntimeFlavor|TestRunArgsUsesConfiguredBinaryName -v
```

Expected: PASS.

- [ ] **Step 7: Build both binaries**

Run:

```powershell
go build ./cmd/gentle-ai
go build ./cmd/gentle-ai-crisp
```

Expected: both builds succeed. This task only changes binary identity/help metadata; deeper config-root separation remains guarded by explicit future work before real installation use.

- [ ] **Step 8: Commit this task**

```powershell
git add cmd/gentle-ai-crisp internal/app/flavor.go internal/app/flavor_test.go internal/app/app.go
git commit -m "feat: add experimental crisp binary flavor"
```

---

## Task 8: Add design traceability and full regression check

**Files:**
- Modify: `docs/superpowers/specs/2026-05-13-sdd-crisp-profile-design.md` only if implementation changed a design decision.

- [ ] **Step 1: Run focused SDD component tests**

Run:

```powershell
go test ./internal/components/sdd -v
```

Expected: PASS.

- [ ] **Step 2: Run full Go test suite**

Run:

```powershell
go test ./...
```

Expected: PASS. If Windows symlink permissions cause known failures, capture the exact failing test and confirm whether it matches the documented Developer Mode/Admin gotcha.

- [ ] **Step 3: Run build checks**

Run:

```powershell
go build ./cmd/gentle-ai
go build ./cmd/gentle-ai-crisp
```

Expected: PASS.

- [ ] **Step 4: Run Python demo tests when Python dependencies are available**

Run:

```powershell
python -m pytest -q testdata/crisp-forecasting-demo
```

Expected: `3 passed`. If dependencies are missing, record the missing package error as an environment limitation and keep the Go fixture tests as the required repo-level validation.

- [ ] **Step 5: Record implementation scope in the design doc**

If the implementation chooses OpenCode-only MVP wiring, append this exact note to the design doc under `## Open decisions for implementation planning`:

```markdown
- MVP implementation wires OpenCode first because the existing SDD profile-overlay path is strongest there. Other adapters should port CRISP assets after OpenCode parity is validated.
```

If implementation wires more than OpenCode, do not add the note.

- [ ] **Step 6: Commit this task**

```powershell
git add docs/superpowers/specs/2026-05-13-sdd-crisp-profile-design.md
git commit -m "docs: update crisp implementation traceability"
```

If the design doc did not change, skip the commit and record “no docs changes” in the final implementation summary.

---

## Self-review checklist

### Spec coverage

- Native SDD-CRISP profile: Tasks 1, 2, 5.
- CRISP phase agents: Tasks 1, 3.
- Result contract and Engram protocol: Tasks 1, 3, 4.
- Python/ML adaptive direction: Task 6 fixture plus CRISP skill rules; deeper repo detection remains a follow-up after prompt/profile parity.
- Governance progressive guardrails: Task 3 skill rules and Task 4 parity markers.
- Separate binary/config/profile strategy: Task 7 establishes separate binary flavor; full config-root isolation is outside this MVP and must be designed before real shared-config installs.
- MVP forecasting demo: Task 6.
- Structural/operational parity: Task 4.

### Known implementation boundaries

- This plan wires OpenCode first because existing profile overlay support is strongest there.
- This plan does not claim full config-root isolation for every adapter. It creates the experimental binary identity and keeps shared-config writes explicit.
- This plan does not require live GitLab, Jenkins, GCP, Postgres, APIs, or MLflow.

### Verification commands

```powershell
go test ./internal/components/sdd -v
go test ./internal/app -v
go test ./...
go build ./cmd/gentle-ai
go build ./cmd/gentle-ai-crisp
python -m pytest -q testdata/crisp-forecasting-demo
```

### Review workload forecast

- Estimated changed lines: high, likely over 400 lines if implemented as one PR.
- Chained PRs recommended: Yes.
- Suggested slices:
  1. CRISP prompt/profile catalog and parity tests.
  2. CRISP skill assets and orchestrator addendum.
  3. OpenCode injection wiring.
  4. Forecasting demo fixture.
  5. Experimental binary flavor.
