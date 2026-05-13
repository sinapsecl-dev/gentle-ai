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

func TestRunArgsUsesConfiguredBinaryNameForUnknownCommandHint(t *testing.T) {
	original := RuntimeFlavor
	RuntimeFlavor = Flavor{Name: "gentle-ai-crisp", ConfigNamespace: "gentle-ai-crisp"}
	t.Cleanup(func() { RuntimeFlavor = original })

	err := RunArgs([]string{"definitely-unknown-command"}, &bytes.Buffer{})
	if err == nil {
		t.Fatal("RunArgs() error = nil, want unknown command error")
	}
	if !strings.Contains(err.Error(), "gentle-ai-crisp help") {
		t.Fatalf("unknown command error = %q, want crisp help hint", err.Error())
	}
}

func TestRunArgsUsesConfiguredBinaryNameForHelp(t *testing.T) {
	original := RuntimeFlavor
	RuntimeFlavor = Flavor{Name: "gentle-ai-crisp", ConfigNamespace: "gentle-ai-crisp"}
	t.Cleanup(func() { RuntimeFlavor = original })

	var out bytes.Buffer
	if err := RunArgs([]string{"help"}, &out); err != nil {
		t.Fatalf("RunArgs(help) error = %v", err)
	}
	if !strings.Contains(out.String(), "gentle-ai-crisp") {
		t.Fatalf("help output = %q, want gentle-ai-crisp", out.String())
	}
}
