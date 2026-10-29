<!-- ![Use Package logo](UseLogo) -->

# `enum` Intro

Provides and reduces code lines for dealing with common and general enums + detects specific os distros correctly.

## Git Clone

`git clone https://gitlab.com/auk-go/enum.git`

### Prerequisites

- Update git to latest 2.29
- Update or install the latest of Go 1.17.8
- Either add your ssh key to your gitlab account
- Or, use your access token to clone it.

## Installation

`go get gitlab.com/auk-go/enum`

## Why enum?

Provides and reduces code lines for dealing with common and general enums + detects specific os distros correctly.

## Examples

### Check OS

```
osmixtype.Ubuntu.IsMajorAtLeast(18) :  false
osmixtype.Ubuntu.IsMajorAtLeast(20) :  false
osmixtype.Ubuntu.IsMajorAtLeast(21) :  false
osmixtype.Windows.IsMajorAtLeast(10) :  true
osmixtype.Windows.IsWindows11() :  true
osmixtype.Windows.Name():  Windows
osmixtype.Windows.ProductName():  Windows 11.22621 Professional
osmixtype.Windows.RawProductName():  Windows 10 Pro
```

### Check OS - Ubuntu 18

```
osmixtype.Ubuntu.IsMajorAtLeast(18) :  true
```

### Check OS - Debian

```
osmixtype.Debian.IsMajorAtLeast(5) :  true
```

### Check OS - CentOS

```
osmixtype.Centos.IsMajorAtLeast(7) :  true
```

### Check OS - MacOS

```
osmixtype.ma.IsMajorAtLeast(10) :  true
```

### Brackets

```go
bracket := brackets.Parenthesis

fmt.Println(bracket.WrapAny("something")) // (something)

bracket2 := brackets.ParenthesisStart // (
fmt.Printf("sizeof(bracket2) = %d\n", unsafe.Sizeof(bracket2)) // 1

someTrue := true
Val := issetter.True

fmt.Printf("sizeof(someTrue) = %d\n", unsafe.Sizeof(&someTrue)) // 8
fmt.Printf("sizeof(Val) = %d\n", unsafe.Sizeof(Val)) // 8

fmt.Println(bracket2.WrapAny("something2")) // (something2)
fmt.Println(bracket2.WrapFmtString("something to do with {wrapped}", "something2")) // replaces '{wrapped}' with (something2)
fmt.Println(bracket2.WrapSkipOnExist("(something2)")) // (something2)
fmt.Println(bracket2.IsWrapped("(something2)")) // true
```

Similar goes for quotes...

### String Types Test

```go
fmt.Println(instructiontype.New("DependsOnx")) // Invalid typename:[instructiontype.Variant], value given : ["DependsOnx"], cannot find in the enum map. reference values can be [Invalid[0], Scoping[1], DependsOn[2], InstallPackages[3], OsServices[4], EnvironmentVariables[5], EnvironmentPaths[6], Ssl[7], CronTabs[8], Ethernet[9], PowerDns[10], MySql[11], PostgreSql[12], PhpMyAdmin[13], PhpPgAdmin[14], Nginx[15], Apache[16], PureFtp[17], FtpUser[18], Compress[19], Decompress[20], DownloadDecompress[21], OsUsersManage[22], OsGroupsManage[23], OsServicesCreate[24], LinuxServiceCreate[25], Ssh[26], PathModifiers[27], FileSystem[28], ScriptsCollection[29], Firewall[30], FirewallIpTables[31], CliInstructions[32], InstructionsExecuteAfterReboot[33], Verify[34]], brackets values can be used to unmarshal as well.
fmt.Println(dbaction.New("Create")) // Create <nil> -- no error

alimStrType := strtype.New("alimx") // alimx -- created no print

fmt.Println(alimStrType.SafeSubString(0, 1)) // a
fmt.Println(alimStrType.SafeSubString(0, alimStrType.Length())) // alimx
fmt.Println(alimStrType.SafeSubString(-1, alimStrType.Length()+5)) // alimx
fmt.Println(alimStrType.SafeSubStringStart(2)) // imx
fmt.Println(alimStrType.SafeSubStringEnd(2)) // al
fmt.Println(alimStrType.SafeSplit(2)) // al imx
fmt.Println(alimStrType.SafeSplit(alimStrType.Length() + 5)) // alimx
```

## Contributors

- [Alim Ul Karim](https://www.google.com/search?q=Alim+Ul+Karim) 

## License

[MIT License](/LICENSE)

## Memory

- .ai-memory/memory/readme.md
- .ai-memory/memory/learned/01-recent-session-2026-10-29.md
- .ai-memory/plans/readme.md
