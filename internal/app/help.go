package app

import (
	"fmt"
	"io"
)

func printHelp(w io.Writer, binaryName, version string) {
	fmt.Fprintf(w, `%s — Gentle-AI: Ecosystem, Frameworks, Workflows (%s)

USAGE
  %s                     Launch interactive TUI
  %s <command> [flags]

COMMANDS
  install      Configure AI coding agents on this machine
  uninstall    Remove Gentle AI managed files from this machine
  sync         Sync agent configs and skills to current version
  skill-registry refresh
               Refresh .atl/skill-registry.md with cache-hit fast path
  update       Check for available updates
  upgrade      Apply updates to managed tools
  restore      Restore a config backup
  version      Print version

FLAGS
  --help, -h    Show this help

Run '%s help' for this message.
Documentation: https://github.com/Gentleman-Programming/gentle-ai
`, binaryName, version, binaryName, binaryName, binaryName)
}
