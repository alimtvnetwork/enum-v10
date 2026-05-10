// AU (Cycle 53) — §09 row-172 advisory contract pin.
//
// Row 172 in spec/01-app/09-converters.md ("Using `*WithDefault` then
// re-validating hides the malformed input") was reclassified ⓘ-advisory
// in Cycle 49 because the claim is dev-experience, not a code property.
// AU pins the *underlying* code-level invariant that makes the advisory
// true: the boolean second-return is the SOLE discriminator between
// "input parsed and happened to equal default" vs "input was malformed
// and was substituted with default". If a future upstream change ever
// returned `true` on the substitution branch (or `false` on a successful
// parse that yields default), this test fails — and the advisory's
// premise has silently collapsed.
//
// Mirrors the AJ-46 pattern: contract tests in tests/contracttests/
// against `core-v9 v1.5.8` upstream `converters` package.
package contracttests

import (
	"strconv"
	"testing"

	"github.com/alimtvnetwork/core-v9/converters"
)

// --------------------------------------------------------------------
// Row 172 / IntegerWithDefault — discriminator invariant.
// --------------------------------------------------------------------

func TestConverters_IntegerWithDefault_OkIsSoleDiscriminator(t *testing.T) {
	const def = -1

	// Case A: input parses to a value that *coincides* with the default.
	// The numeric return is indistinguishable from the substitution case;
	// `ok=true` is the ONLY signal that the input was real.
	gotVal, gotOk := converters.StringTo.IntegerWithDefault(strconv.Itoa(def), def)
	if gotVal != def {
		t.Fatalf("legit input parsing to default: got val=%d want %d", gotVal, def)
	}
	if !gotOk {
		t.Fatalf("legit input parsing to default: got ok=false; row-172 advisory premise broken — re-validating the value alone would now wrongly flag a real input as malformed")
	}

	// Case B: malformed input substitutes the default with ok=false.
	// If `ok` were `true` here, a downstream "re-validate" check that
	// only inspects the value would silently accept malformed input as
	// `def` — exactly the failure mode row-172 warns about.
	for _, bad := range []string{"", "abc", "1.2.3", "++1", " 42", "9999999999999999999999"} {
		gotVal, gotOk = converters.StringTo.IntegerWithDefault(bad, def)
		if gotVal != def {
			t.Errorf("bad input %q: substituted val=%d want %d", bad, gotVal, def)
		}
		if gotOk {
			t.Errorf("bad input %q: ok=true (advisory premise broken — substitution would be invisible to value-only re-validation)", bad)
		}
	}
}

// --------------------------------------------------------------------
// Row 172 / ByteWithDefault — same discriminator invariant on the
// byte-typed family. Pinned because §09 lists `ByteWithDefault` in the
// same advisory cluster.
// --------------------------------------------------------------------

func TestConverters_ByteWithDefault_OkIsSoleDiscriminator(t *testing.T) {
	const def byte = 7

	// Case A: legit input that equals default.
	gotVal, gotOk := converters.StringTo.ByteWithDefault(strconv.Itoa(int(def)), def)
	if gotVal != def {
		t.Fatalf("legit byte parsing to default: got %d want %d", gotVal, def)
	}
	if !gotOk {
		t.Fatalf("legit byte parsing to default: ok=false; row-172 byte-family premise broken")
	}

	// Case B: malformed / out-of-range inputs substitute default with ok=false.
	for _, bad := range []string{"", "abc", "-1", "256", "1.5"} {
		gotVal, gotOk = converters.StringTo.ByteWithDefault(bad, def)
		if gotVal != def {
			t.Errorf("bad byte input %q: substituted val=%d want %d", bad, gotVal, def)
		}
		if gotOk {
			t.Errorf("bad byte input %q: ok=true (advisory premise broken)", bad)
		}
	}
}

// --------------------------------------------------------------------
// Row 172 — re-validation hazard scenario.
// Demonstrates the exact developer mistake the advisory warns about,
// using only the public API. If the assertion below ever flips, it
// means the substitution branch became distinguishable by value alone
// (e.g. via a sentinel out-of-band marker) — at which point row 172
// should be promoted from ⓘ-advisory to ✅ with an updated rationale.
// --------------------------------------------------------------------

func TestConverters_IntegerWithDefault_ValueOnlyRevalidationHidesMalformed(t *testing.T) {
	const def = 0 // a common "zero-value default" that is also a legit parse target

	legitVal, _ := converters.StringTo.IntegerWithDefault("0", def)
	subVal, _ := converters.StringTo.IntegerWithDefault("not-a-number", def)

	// The hazard: a developer who re-validates by checking the returned
	// value against `def` cannot tell these two cases apart.
	if legitVal != subVal {
		t.Fatalf("row-172 hazard scenario broken: legit=%d substituted=%d — values now differ, advisory should be promoted to ✅", legitVal, subVal)
	}
}
