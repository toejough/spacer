//go:build targ

package dev

import "github.com/toejough/targ"

func init() {
	targ.Register(
		targ.Targ("npx vite").Name("dev"),
		targ.Targ("npx vitest run").Name("test"),
		targ.Targ("npx vue-tsc --noEmit && npx vite build").Name("build"),
		targ.Targ("npx vue-tsc --noEmit && npx vitest run").Name("check"),
		targ.Targ("rm -rf dist/ test-results/").Name("clean"),
	)
}
