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
	"crisp-business":           "Define business objective, hypothesis, stakeholder, success metric, and risk level",
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
