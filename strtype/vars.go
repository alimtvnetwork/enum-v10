package strtype

import (
	"sync"

	"github.com/alimtvnetwork/core-v9/coredata/coredynamic"
	"github.com/alimtvnetwork/core-v9/coredata/corejson"
	"github.com/alimtvnetwork/core-v9/ostype"
	"github.com/alimtvnetwork/enum-v9/osarchs"
)

var (
	invalidMaps = map[string]bool{
		"":              true,
		"Invalid":       true,
		"invalid":       true,
		"Unknown":       true,
		"unknown":       true,
		"Undefined":     true,
		"undefined":     true,
		"Uninitialized": true,
		"uninitialized": true,
		"None":          true,
		"none":          true,
		"Unspecified":   true,
		"unspecified":   true,
	}

	globalMutex = sync.Mutex{}

	// Arch
	// Current OS architecture
	Arch  = osarchs.CurrentArch
	Group = ostype.CurrentGroupVariant.Group
	// Type Current Os Type
	Type               = ostype.CurrentGroupVariant
	bytesSerializer    = corejson.Serialize.ToBytesErr
	stringDeserializer = corejson.Deserialize.BytesTo.String

	typeName = coredynamic.SafeTypeName(Variant(""))
)
