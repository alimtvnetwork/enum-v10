package osdetect

import (
	"path/filepath"
	
	"github.com/alimtvnetwork/core-v9/chmodhelper"
	"github.com/alimtvnetwork/core-v9/coredata/coreonce"
	"github.com/alimtvnetwork/core-v9/coreimpl/enumimpl"
	"github.com/alimtvnetwork/core-v9/errcore"
)

var (
	Ranges = [...]string{
		Invalid:               "Invalid",
		AnyOs:                 "AnyOs",
		Windows:               "Windows",
		Unix:                  "Unix",
		Linux:                 "Linux",
		MacOs:                 "MacOs",
		Ubuntu:                "Ubuntu",
		Debian:                "Debian",
		ArchLinux:             "ArchLinux",
		FreeBsd:               "FreeBsd",
		Centos:                "Centos",
		RedHatEnterpriseLinux: "RedHatEnterpriseLinux",
		Docker:                "Docker",
		Android:               "Android",
	}
	
	lowerCaseNames = [...]string{
		Invalid:               "invalid",
		AnyOs:                 "any-os",
		Windows:               "windows",
		Unix:                  "unix",
		Linux:                 "linux",
		MacOs:                 "macOs",
		Ubuntu:                "ubuntu",
		Debian:                "Debian",
		ArchLinux:             "ArchLinux",
		FreeBsd:               "FreeBsd",
		Centos:                "centos",
		RedHatEnterpriseLinux: "redhat-enterprise-linux",
		Docker:                "docker",
		Android:               "android",
	}
	
	osGroupMap = map[Variant]bool{
		Windows: true,
		Unix:    true,
		Android: true,
	}
	
	linuxGroupMap = map[Variant]bool{
		Linux:  true,
		Ubuntu: true,
		Centos: true,
	}
	
	aliasMap = map[string]byte{
		"all":     AnyOs.Value(),
		"any-os":  AnyOs.Value(),
		"anyOs":   AnyOs.Value(),
		"all-os":  AnyOs.Value(),
		"default": AnyOs.Value(),
		"Default": AnyOs.Value(),
		"win":     Windows.Value(),
		"windows": Windows.Value(),
		"unix":    Unix.Value(),
		macOsName: MacOs.Value(),
		"mac":     MacOs.Value(),
		"macOs":   MacOs.Value(),
		"macOS":   MacOs.Value(),
		"dar":     MacOs.Value(),
		"darwin":  MacOs.Value(),
		"Darwin":  MacOs.Value(),
		"linux":   Linux.Value(),
	}
	
	// winRegistryKeyNames
	//
	// Registry location : "Computer\HKEY_LOCAL_MACHINE\SOFTWARE\Microsoft\Windows NT\CurrentVersion"
	// https://prnt.sc/Yf7QKMuJNTse
	//
	// go - How to return a default value from windows/registry with golang - Stack Overflow](https://stackoverflow.com/questions/36998532/how-to-return-a-default-value-from-windows-registry-with-golang)
	// Detect Windows version : https://stackoverflow.com/questions/44363911/detect-windows-version-in-go
	// Operating System Version : https://docs.microsoft.com/en-us/windows/win32/sysinfo/operating-system-version?redirectedfrom=MSDN
	//
	// Issue : https://github.com/alimtvnetwork/enum-v10/-/issues/4
	winRegistryKeyNames = windowsCurrentVersionRegistry{
		buildBranch:          "BuildBranch",
		productName:          "ProductName",
		installationType:     "InstallationType", // eg. Client (Win 10, 11) or Server (2016, 2019)
		compositionEditionId: "CompositionEditionID",
		editionId:            "EditionID",
		releaseId:            "ReleaseId", // eg. 2009, 1809
		pathName:             "SystemRoot",
		systemRoot:           "SystemRoot",
		currentVersion:       "CurrentVersion",
		majorVersionNumber:   "CurrentMajorVersionNumber",
		minorVersionNumber:   "CurrentMinorVersionNumber",
		currentBuild:         "CurrentBuild",
		currentBuildNumber:   "CurrentBuildNumber",
		registeredOwner:      "RegisteredOwner",
	}
	
	generateInstance             = generate{}
	currentOsDetailGeneratorOnce = coreonce.NewAnyOnce(func() interface{} {
		osDetail, err := generateInstance.OperatingSystemDetailLazy()
		
		return &OsDetailWithErr{
			OperatingSystemDetail: osDetail,
			Error:                 errcore.ToString(err),
		}
	})
	
	currentOsMixTypeOnce = coreonce.NewByteOnce(func() byte {
		osDetail, err := GetCurrentOsDetail()
		
		if err != nil {
			return Invalid.Value()
		}
		
		return osDetail.OsMixType.Value()
	})
	
	currentOsMixTypesMapOnce = coreonce.NewAnyOnce(func() interface{} {
		return generateInstance.currentOsMixTypesMap()
	})
	
	currentOsMixTypesOnce = coreonce.NewAnyOnce(func() interface{} {
		return generateInstance.currentOsMixTypes()
	})
	
	// path example : /var/temp/os-detail or
	//                c:\windows\temp\os-detail (windows)
	osDetailTempCacheRootPath = filepath.Join(
		chmodhelper.TempDirGetter.TempPermanent(),
		osDetailTempDirName,
	)
	
	osDetailTempCachePath = filepath.Join(
		osDetailTempCacheRootPath,
		osCachedTempFileName)
	
	BasicEnumImpl = enumimpl.New.BasicByte.UsingFirstItemSliceAliasMap(
		Invalid,
		Ranges[:],
		aliasMap)
)
