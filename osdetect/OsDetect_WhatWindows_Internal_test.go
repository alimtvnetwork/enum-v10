package osdetect

// AQ (Cycle 53): white-box coverage for the unexported
// `WindowsSystemDetail.initialize` / `whatWindows` branch matrix.
// These methods build the canonical Windows display name and are
// 100% data-driven (no registry / syscall dependency), so they can
// run on any Linux CI runner. Prior to this file every branch
// inside `whatWindows` was at 0% because the public uplift suite
// in `osdetect_test` package can't reach unexported methods.

import (
	"strings"
	"testing"

	"github.com/alimtvnetwork/enum-v10/inttype"
	"github.com/alimtvnetwork/enum-v10/strtype"
)

// helper — keeps each subtest body small and intent-revealing.
func wsdRequire(t *testing.T, label, got, wantSubstr string) {
	t.Helper()
	if !strings.Contains(got, wantSubstr) {
		t.Errorf("%s: %q does not contain %q", label, got, wantSubstr)
	}
}

func Test_WindowsSystemDetail_WhatWindows_AllBranches(t *testing.T) {
	t.Run("nil receiver returns empty", func(t *testing.T) {
		var nilWsd *WindowsSystemDetail
		if got := nilWsd.whatWindows("ignored"); got != "" {
			t.Errorf("nil receiver should yield empty Variant, got %q", got)
		}
	})

	t.Run("server branch wins over client checks", func(t *testing.T) {
		wsd := &WindowsSystemDetail{
			IsServer:       true,
			ServerVersion:  inttype.Variant(2019),
			Edition:        strtype.Variant("ServerStandard"),
			CurrentBuildId: inttype.Variant(17763),
		}
		got := string(wsd.whatWindows("ignored"))
		wsdRequire(t, "server", got, "Windows Server")
		wsdRequire(t, "server", got, "ServerStandard")
	})

	t.Run("windows 11 branch (CurrentBuildId >= 22000)", func(t *testing.T) {
		wsd := &WindowsSystemDetail{
			IsClient:        true,
			WindowsVersion:  inttype.Variant(11),
			CurrentBuildId:  inttype.Variant(windows11BuildIdentifier + 1),
			Edition:         strtype.Variant("Professional"),
		}
		got := string(wsd.whatWindows("ignored"))
		wsdRequire(t, "win11", got, "Windows 11.")
		wsdRequire(t, "win11", got, "Professional")
	})

	t.Run("windows 10 branch", func(t *testing.T) {
		wsd := &WindowsSystemDetail{
			IsClient:        true,
			WindowsVersion:  inttype.Variant(10),
			CurrentBuildId:  inttype.Variant(19045),
			Edition:         strtype.Variant("Enterprise"),
		}
		got := string(wsd.whatWindows("ignored"))
		wsdRequire(t, "win10", got, "Windows 10.")
		wsdRequire(t, "win10", got, "Enterprise")
	})

	t.Run("windows 8 branch", func(t *testing.T) {
		wsd := &WindowsSystemDetail{
			IsClient:        true,
			WindowsVersion:  inttype.Variant(8),
			CurrentBuildId:  inttype.Variant(9600),
			Edition:         strtype.Variant("Workstation"),
		}
		got := string(wsd.whatWindows("ignored"))
		wsdRequire(t, "win8", got, "Windows 8.")
		wsdRequire(t, "win8", got, "Workstation")
	})

	t.Run("fallback branch returns provided windowsName", func(t *testing.T) {
		// Not server, not 8/10/11 — the function returns the raw
		// `windowsName` argument verbatim.
		wsd := &WindowsSystemDetail{
			IsClient:        true,
			WindowsVersion:  inttype.Variant(7), // Windows 7 — no dedicated branch
			CurrentBuildId:  inttype.Variant(7601),
		}
		if got := string(wsd.whatWindows("Windows 7 Ultimate")); got != "Windows 7 Ultimate" {
			t.Errorf("fallback should return the raw windowsName, got %q", got)
		}
	})
}

func Test_WindowsSystemDetail_Initialize_PopulatesGeneratedName(t *testing.T) {
	wsd := &WindowsSystemDetail{
		IsClient:       true,
		WindowsVersion: inttype.Variant(10),
		CurrentBuildId: inttype.Variant(19045),
		Edition:        strtype.Variant("Enterprise"),
	}

	// initialize is unexported; we drive it via white-box access.
	wsd.initialize("ignored-by-win10-branch")

	got := string(wsd.GeneratedWindowsName)
	if !strings.Contains(got, "Windows 10.") || !strings.Contains(got, "Enterprise") {
		t.Errorf("initialize should set GeneratedWindowsName to a Windows 10 string, got %q", got)
	}
}

func Test_WindowsSystemDetail_Initialize_FallbackUsesProvidedName(t *testing.T) {
	wsd := &WindowsSystemDetail{
		IsClient:       true,
		WindowsVersion: inttype.Variant(7), // no Win7 branch -> fallback
	}
	wsd.initialize("Custom Windows Build")
	if string(wsd.GeneratedWindowsName) != "Custom Windows Build" {
		t.Errorf("initialize fallback should pass windowsName through, got %q", wsd.GeneratedWindowsName)
	}
}
