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
	for _, sddPhase := range []string{"sdd-spec", "sdd-design", "sdd-tasks", "sdd-apply", "sdd-verify"} {
		taskPerms[sddPhase] = "allow"
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
