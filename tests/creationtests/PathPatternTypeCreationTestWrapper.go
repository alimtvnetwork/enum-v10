package creationtests

import "github.com/alimtvnetwork/enum-v9/pathpatterntype"

type PathPatternTypeCreationTestWrapper struct {
	PathType                      pathpatterntype.Variant
	Name, FullName, CurlyFullName string
	CompiledTemplateFullPath      string
	AssociatedTemplatePaths       []string
}
