package scripttype

import (
	"fmt"
	"strings"
)

type ScriptDefault struct {
	ScriptType       Variant
	IsImplemented    bool
	ProcessName      string
	DefaultArguments []string
}

// String returns a deterministic textual representation of the script default.
// Per RCA Pattern P9 (.lovable/memory/07-test-failure-rca-patterns.md) we
// MUST NOT delegate to `converters.AnyTo.ValueString(*it)` (or any %v sprint)
// — formatting fields explicitly is recursion-safe regardless of receiver
// shape changes.
func (it *ScriptDefault) String() string {
	if it == nil {
		return "<nil>"
	}
	return fmt.Sprintf(
		"{ScriptType:%s IsImplemented:%t ProcessName:%s DefaultArguments:[%s]}",
		it.ScriptType.String(),
		it.IsImplemented,
		it.ProcessName,
		strings.Join(it.DefaultArguments, " "),
	)
}

