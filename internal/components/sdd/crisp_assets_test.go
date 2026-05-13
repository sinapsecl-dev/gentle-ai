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
