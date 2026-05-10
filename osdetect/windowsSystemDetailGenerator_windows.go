package osdetect

import (
	"github.com/alimtvnetwork/core-v9/codestack"
	"github.com/alimtvnetwork/core-v9/converters"
	"github.com/alimtvnetwork/core-v9/coreversion"
	"github.com/alimtvnetwork/core-v9/errcore"
	"github.com/alimtvnetwork/core-v9/ostype"
	"github.com/alimtvnetwork/enum-v9/inttype"
	"github.com/alimtvnetwork/enum-v9/osarchs"
	"github.com/alimtvnetwork/enum-v9/strtype"
	"golang.org/x/sys/windows/registry"
)

type windowsSystemDetailGenerator struct {
	rawErrCollection errcore.RawErrCollection
	rootRegistryKey  registry.Key
}

func (it *windowsSystemDetailGenerator) Value(
	name string,
) strtype.Variant {
	currentValue, _, err := it.
		rootRegistryKey.
		GetStringValue(name)
	
	if err != nil {
		it.rawErrCollection.
			AddWithRef(err, name)
	}
	
	return strtype.New(currentValue)
}

func (it *windowsSystemDetailGenerator) ValueInt(
	name string,
) inttype.Variant {
	currentValue, _, err := it.
		rootRegistryKey.
		GetIntegerValue(name)
	
	if err != nil {
		it.rawErrCollection.
			AddWithRef(err, name)
	}
	
	return inttype.Variant(currentValue)
}

func (it *windowsSystemDetailGenerator) CloseRegKeyRead() {
	err := it.rootRegistryKey.Close()
	
	it.rawErrCollection.AddError(err)
}

func (it *windowsSystemDetailGenerator) Finalize() error {
	it.CloseRegKeyRead()
	
	return it.rawErrCollection.
		CompiledErrorWithStackTraces()
}

func (it windowsSystemDetailGenerator) CompiledErrorWithStackTraces() error {
	if it.rawErrCollection.IsEmpty() {
		return nil
	}
	
	stackTraces := codestack.StacksTo.StringDefault()
	it.rawErrCollection.AddString(stackTraces)
	
	return it.rawErrCollection.CompiledErrorWithStackTraces()
}

// SystemDetail
//
//	Reference : https://github.com/alimtvnetwork/enum-v9/-/issues/4
func (it windowsSystemDetailGenerator) SystemDetail() (*OperatingSystemDetail, error) {
	buildBranch := it.Value(winRegistryKeyNames.buildBranch)
	productName := it.Value(winRegistryKeyNames.productName)
	installationType := it.Value(winRegistryKeyNames.installationType)         // client or server
	editionId := it.Value(winRegistryKeyNames.editionId)                       // ServerStandard or Profession or Workstation
	compositionEditionId := it.Value(winRegistryKeyNames.compositionEditionId) // ServerStandard or Profession or Workstation
	releaseId := it.Value(winRegistryKeyNames.releaseId)
	systemRoot := it.Value(winRegistryKeyNames.systemRoot)
	currentVersionValue := it.Value(winRegistryKeyNames.currentVersion) // example CurrentVersion: 17763 (not right)
	majorVersion := it.ValueInt(winRegistryKeyNames.majorVersionNumber)
	minorVersion := it.ValueInt(winRegistryKeyNames.minorVersionNumber)
	buildNumber := it.Value(winRegistryKeyNames.currentBuildNumber)
	registerOwner := it.Value(winRegistryKeyNames.registeredOwner)
	
	currentVersion := coreversion.New.MajorMinorBuild(
		majorVersion.String(),
		minorVersion.String(),
		currentVersionValue.String())
	
	compiledVersion := coreversion.New.MajorMinorBuild(
		majorVersion.String(),
		minorVersion.String(),
		buildNumber.String())
	
	versionNumberStrings := windowsVersionNumberLazyRegex.
		CompileMust().
		FindStringSubmatch(productName.String()) // From "Windows Server 2019 Standard" to 2019
	
	hasVersionNumber := len(versionNumberStrings) > 0
	
	windowsVersion := inttype.Zero
	serverVersion := inttype.Zero
	isServer := installationType == windowsServerInstallationType
	var versionInNumber int
	
	if hasVersionNumber {
		toNumber, err := converters.StringTo.Integer(
			versionNumberStrings[0])
		
		it.rawErrCollection.AddError(err)
		versionInNumber = toNumber
	}
	
	if isServer {
		serverVersion = inttype.Variant(versionInNumber)
	} else {
		windowsVersion = inttype.Variant(versionInNumber)
	}
	
	finalErr := it.Finalize()
	
	winDetail := WindowsSystemDetail{
		WindowsVersion:     windowsVersion,
		ServerVersion:      serverVersion,
		CurrentVersion:     currentVersion.NonPtr(),
		CompiledVersion:    compiledVersion.NonPtr(),
		ReleaseId:          releaseId.IntType(),
		CurrentBuildId:     buildNumber.IntType(),
		BuildBranch:        buildBranch,
		InstallType:        installationType,
		SystemRoot:         systemRoot,
		Edition:            editionId,
		CompositionEdition: compositionEditionId,
		RegisteredOwner:    registerOwner,
		IsServer:           isServer,
		IsClient:           installationType == windowsClientInstallationType,
		CompiledError:      finalErr,
	}
	
	winDetail.initialize(productName.String())
	
	osDetail := OperatingSystemDetail{
		Name:          winDetail.GeneratedWindowsName,
		ProductName:   productName,
		Vendor:        windowsVendor,
		Version:       strtype.New(compiledVersion.String()),
		Release:       releaseId,
		Architecture:  osarchs.CurrentArch,
		OsMixType:     Windows,
		WindowsDetail: &winDetail,
		Group:         ostype.WindowsGroup,
		IsDocker:      IsRunningInDockerContainer(),
	}
	
	return &osDetail, finalErr
}

func (it windowsSystemDetailGenerator) AsWindowsSysDetailDefiner() windowsSysDetailDefiner {
	return &it
}
