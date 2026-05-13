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
	for _, sddPhase := range []string{"sdd-spec", "sdd-design", "sdd-tasks", "sdd-apply", "sdd-verify"} {
		if taskMap[sddPhase] != "allow" {
			t.Fatalf("permission for %s = %v, want allow", sddPhase, taskMap[sddPhase])
		}
	}
}

func TestGenerateCRISPProfileOverlayPromptContainsSDDHandoffContract(t *testing.T) {
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
	prompt := orch["prompt"].(string)

	markers := []string{
		"must NOT implement product code inline",
		"planning/governance/review",
		"CRISP-to-SDD handoff",
		"business-understanding",
		"data-understanding",
		"data-preparation",
		"modeling",
		"evaluation",
		"deployment",
		"binding constraints",
		"sdd-spec",
		"sdd-design",
		"sdd-tasks",
		"sdd-apply",
		"sdd-verify",
		"missing CRISP artifacts",
		"ask for confirmation",
	}
	for _, marker := range markers {
		if !strings.Contains(prompt, marker) {
			t.Fatalf("orchestrator prompt missing marker %q", marker)
		}
	}
}
