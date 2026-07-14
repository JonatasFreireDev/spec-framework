package specframework

import "embed"

// Assets contains the product starter and framework-owned runtime files.
// Init composes product-owned starter assets through versioned contracts; method
// assets are materialized in the user cache so the released CLI works offline.
//
//go:embed FRAMEWORK.md all:starter all:examples/events docs/artifact-registry-modules.md docs/execution-runtime.md docs/engineering-systems.md docs/lifecycle-and-approvals.md framework/AGENTS.framework.md framework/delivery-closure.md framework/init framework/skills framework/template
var Assets embed.FS
