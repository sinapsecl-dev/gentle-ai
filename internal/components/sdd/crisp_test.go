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
