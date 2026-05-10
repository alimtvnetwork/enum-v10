package osdetect

import (
	"fmt"
	
	"github.com/alimtvnetwork/core-v9/coreversion"
	"github.com/alimtvnetwork/enum-v9/inttype"
	"github.com/alimtvnetwork/enum-v9/strtype"
)

type WindowsSystemDetail struct {
	WindowsVersion       inttype.Variant     // eg. 8, 10, 11 https://stackoverflow.com/a/69922526
	ServerVersion        inttype.Variant     // eg. 2016, 2019
	CurrentVersion       coreversion.Version // eg. 6.3 https://t.ly/XLFC
	CompiledVersion      coreversion.Version
	ReleaseId            inttype.Variant
	CurrentBuildId       inttype.Variant
	BuildBranch          strtype.Variant
	InstallType          strtype.Variant
	SystemRoot           strtype.Variant
	Edition              strtype.Variant // Example: "ServerStandard", "Professional", "Enterprise", "Workstation"
	CompositionEdition   strtype.Variant // eg. "Enterprise" - For Windows 10, "ServerStandard" -- Windows Server
	GeneratedWindowsName strtype.Variant `json:"WindowsName,omitempty"`
	RegisteredOwner      strtype.Variant `json:"RegisteredOwner,omitempty"` // eg. owner name or email.
	IsServer             bool            // Refers to Windows Server
	IsClient             bool            // Refers to Windows 10
	CompiledError        error
}

func (it *WindowsSystemDetail) IsNull() bool {
	return it == nil
}

func (it *WindowsSystemDetail) WinVer() inttype.Variant {
	if it.IsClient {
		return it.WindowsVersion
	}
	
	if it.IsServer {
		return it.ServerVersion
	}
	
	return inttype.Zero
}

func (it *WindowsSystemDetail) IsDefined() bool {
	return it != nil
}

func (it *WindowsSystemDetail) IsNullOr(isCondition bool) bool {
	return it == nil || isCondition
}

func (it *WindowsSystemDetail) IsDefinedPlus(isCondition bool) bool {
	return it != nil && isCondition
}

func (it WindowsSystemDetail) IsWindows8() bool {
	return it.IsWindowsEqual(8)
}

func (it WindowsSystemDetail) IsWindows11OrAbove() bool {
	return it.IsWindows11()
}

func (it WindowsSystemDetail) IsWindowsGreaterEqual(number int) bool {
	if it.IsNullOr(it.IsServer) {
		return false
	}
	
	return it.WindowsVersion.IsGreaterEqual(number)
}

func (it *WindowsSystemDetail) IsWindowsServerGreaterEqual(number int) bool {
	if it.IsNullOr(it.IsClient) {
		return false
	}
	
	return it.ServerVersion.IsGreaterEqual(number)
}

func (it WindowsSystemDetail) IsWindowsEqual(number int) bool {
	if it.IsNullOr(it.IsServer) {
		return false
	}
	
	return it.WindowsVersion.IsEqual(number)
}

func (it *WindowsSystemDetail) IsWindowsServerEqual(number int) bool {
	if it.IsNullOr(it.IsClient) {
		return false
	}
	
	return it.ServerVersion.IsEqual(number)
}

func (it WindowsSystemDetail) IsWindows7() bool {
	return it.IsWindowsEqual(7)
}

func (it WindowsSystemDetail) IsWindows10() bool {
	return it.IsWindowsEqual(10)
}

// IsWindows11
//
//	https://t.ly/Jsr1,
//	https://prnt.sc/wAZ5uQScNqk_
func (it WindowsSystemDetail) IsWindows11() bool {
	return it.CurrentBuildId.IsGreaterEqual(
		windows11BuildIdentifier)
}

func (it WindowsSystemDetail) IsWindowsSever() bool {
	return it.IsDefinedPlus(it.IsServer)
}

func (it WindowsSystemDetail) IsWindowsSever2016() bool {
	return it.IsWindowsServerEqual(2016)
}

func (it WindowsSystemDetail) IsWindowsSever2019() bool {
	return it.IsWindowsServerEqual(2019)
}

func (it WindowsSystemDetail) IsWindowsSeverGreaterEqual2016() bool {
	// Bugfix (v0.79.0): previously dispatched to IsWindowsGreaterEqual,
	// which short-circuits to false when IsServer=true — making this
	// helper unconditionally return false for any real server. Route to
	// the server-side helper so server flavors compare correctly.
	return it.IsWindowsServerGreaterEqual(2016)
}

func (it WindowsSystemDetail) IsWindowsSeverGreaterEqual2019() bool {
	return it.IsWindowsServerGreaterEqual(2019)
}

func (it *WindowsSystemDetail) initialize(windowsName string) {
	it.GeneratedWindowsName = it.whatWindows(windowsName)
}

func (it *WindowsSystemDetail) whatWindows(windowsName string) strtype.Variant {
	if it == nil {
		return ""
	}
	
	if it.IsServer {
		output := fmt.Sprintf(
			"Windows Server %s %s %s",
			it.ServerVersion,
			it.Edition,
			it.CurrentBuildId)
		
		return strtype.Variant(output)
	}
	
	if it.IsWindows11() {
		output := fmt.Sprintf(
			"Windows 11.%s %s",
			it.CurrentBuildId,
			it.Edition)
		
		return strtype.Variant(output)
	}
	
	if it.IsWindows10() {
		output := fmt.Sprintf(
			"Windows 10.%s %s",
			it.CurrentBuildId,
			it.Edition)
		
		return strtype.Variant(output)
	}
	
	if it.IsWindows8() {
		output := fmt.Sprintf(
			"Windows 8.%s %s",
			it.CurrentBuildId,
			it.Edition)
		
		return strtype.Variant(output)
	}
	
	return strtype.Variant(windowsName)
}
