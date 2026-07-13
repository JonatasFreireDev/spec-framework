package specframework

import "embed"

// Assets contains the product starter and framework-owned runtime files.
// Init copies only starter/product into adopter repositories; method assets are
// materialized in the versioned user cache so the released CLI works offline.
//
//go:embed FRAMEWORK.md all:starter all:examples/events framework/AGENTS.framework.md framework/delivery-closure.md framework/decisions framework/skills framework/template
var Assets embed.FS
