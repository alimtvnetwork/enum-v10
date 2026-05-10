---
name: core-v9 API migration mapping
description: Maps old core-v8 function-level API to core-v9 struct-namespace API discovered via user's PowerShell investigation
type: feature
---

# core-v9 API Migration Mapping

> Discovered 2026-05-06 by user running PowerShell `Select-String` commands against their local `core-v9` checkout (tag v1.5.8).

## Breaking change pattern

`core-v9` moved from **package-level functions** to **struct-based namespaces** exposed as package-level vars in `converters/vars.go`:

```go
var (
    StringsTo     = stringsTo{}
    AnyTo         = anyItemConverter{}
    Map           = convertinternal.Map
    StringTo      = stringTo{}
    PrettyJson    = jsoninternal.Pretty
    JsonString    = jsoninternal.String
    BytesTo       = bytesTo{}
    Integers      = convertinternal.Integers
    KeyValuesTo   = convertinternal.KeyValuesTo
    CodeFormatter = convertinternal.CodeFormatter
)
```

## Confirmed migration map

| Old call (enum-v9 code today) | New call (core-v9 v1.5.8) | Verified? |
|---|---|---|
| `converters.AnyToValueString(x)` | `converters.AnyTo.ValueString(x)` | ✅ Applied 2026-05-06 (6 sites) |
| `converters.StringToInteger(s)` | `converters.StringTo.Integer(s)` | ✅ Applied 2026-05-06 |
| `converters.StringToIntegerWithDefault(s, d)` | `converters.StringTo.IntegerWithDefault(s, d)` | ✅ Applied 2026-05-06 |
| `converters.StringToByte(s)` | `converters.StringTo.Byte(s)` | ✅ Applied 2026-05-06 |
| `converters.StringToIntegerDefault(s)` | `converters.StringTo.IntegerDefault(s)` | ✅ Applied 2026-05-06 |
| `converters.Any.ToFullNameValueString(x)` | `converters.AnyTo.ToFullNameValueString(x)` | ✅ Applied 2026-05-06 (1 site) |
| `coredynamic.TypeName(x)` | `coredynamic.SafeTypeName(x)` | ✅ Applied 2026-05-06 (53 sites) |

## `anyItemConverter` methods (confirmed from user output)

```
ToString(anyItem any) string
String(anyItem any) string
FullString(anyItem any) string
ValueString(anyItem any) string
Bytes(anyItem any) []byte
```

(Plus ~20 more methods; partial list from `Select-String` output.)

## `coredynamic` public functions (confirmed)

- `EmptyAnyCollection()`, `NewAnyCollection(capacity int)`
- `SafeTypeName(item any) string` — **this is the replacement for the missing `TypeName`**
- No public `TypeName` function exists in `coredata/coredynamic/`

## `stringTo` methods — confirmed from spec and applied

The consumer-facing converter spec documents these `core-v9` calls and the
source now compiles against the struct namespace form:

- `converters.StringTo.Integer(s)`
- `converters.StringTo.IntegerWithDefault(s, d)`
- `converters.StringTo.IntegerDefault(s)`
- `converters.StringTo.Byte(s)`

## Where broken calls live in enum-v9

Run this to find all affected files:
```bash
rg -n 'converters\.(AnyToValueString|StringToInteger|StringToByte|StringToIntegerWithDefault|StringToIntegerDefault)' --type go
rg -n 'coredynamic\.TypeName[^s]' --type go
rg -n 'converters\.Any\.' --type go
```

## Next steps

1. Run `./run.ps1 -tc` locally to validate Task AM in the Windows/Go environment.
2. If new undefined symbols appear, inspect `data/coverage/blocked-packages.txt` and update this map.
