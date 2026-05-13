package cli

import (
	"reflect"
	"testing"

	"github.com/gentleman-programming/gentle-ai/internal/model"
)

func TestBuildSDDInjectOptionsEnablesCRISPForCrispFlavorOpenCode(t *testing.T) {
	originalFlavor := RuntimeFlavorName
	RuntimeFlavorName = "gentle-ai-crisp"
	t.Cleanup(func() { RuntimeFlavorName = originalFlavor })

	got := buildSDDInjectOptions(model.Selection{}, "", model.AgentOpenCode)
	want := []model.Profile{{Name: "sdd-crisp"}}
	if !reflect.DeepEqual(got.CRISPProfiles, want) {
		t.Fatalf("CRISPProfiles = %#v, want %#v", got.CRISPProfiles, want)
	}
}

func TestBuildSDDInjectOptionsKeepsDefaultFlavorWithoutCRISP(t *testing.T) {
	originalFlavor := RuntimeFlavorName
	RuntimeFlavorName = "gentle-ai"
	t.Cleanup(func() { RuntimeFlavorName = originalFlavor })

	got := buildSDDInjectOptions(model.Selection{}, "", model.AgentOpenCode)
	if len(got.CRISPProfiles) != 0 {
		t.Fatalf("CRISPProfiles = %#v, want empty", got.CRISPProfiles)
	}
}

func TestBuildSDDInjectOptionsDoesNotEnableCRISPForNonOpenCode(t *testing.T) {
	originalFlavor := RuntimeFlavorName
	RuntimeFlavorName = "gentle-ai-crisp"
	t.Cleanup(func() { RuntimeFlavorName = originalFlavor })

	got := buildSDDInjectOptions(model.Selection{}, "", model.AgentClaudeCode)
	if len(got.CRISPProfiles) != 0 {
		t.Fatalf("CRISPProfiles = %#v, want empty", got.CRISPProfiles)
	}
}
