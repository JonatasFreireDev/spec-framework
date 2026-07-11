package specframework

import "embed"

// Assets contains the product starter and framework-owned installation files.
// Keeping them in the binary makes the released CLI independent from this repository.
//
//go:embed FRAMEWORK.md all:starter framework/AGENTS.framework.md framework/delivery-closure.md framework/decisions framework/skills framework/template
var Assets embed.FS
