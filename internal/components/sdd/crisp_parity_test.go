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
