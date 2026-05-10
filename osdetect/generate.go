package osdetect

import (
	"bufio"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	
	"github.com/alimtvnetwork/core-v9/chmodhelper"
	"github.com/alimtvnetwork/core-v9/constants"
	"github.com/alimtvnetwork/core-v9/coreutils/stringutil"
	"github.com/alimtvnetwork/core-v9/errcore"
	"github.com/alimtvnetwork/core-v9/filemode"
	"github.com/alimtvnetwork/core-v9/osconsts"
	"github.com/alimtvnetwork/core-v9/ostype"
	"github.com/alimtvnetwork/enum-v8/linuxvendortype"
	"github.com/alimtvnetwork/enum-v8/osarchs"
	"github.com/alimtvnetwork/enum-v8/strtype"
)

type generate struct{}

func (it generate) currentOsMixTypes() (results []Variant) {
	results = make(
		[]Variant,
		0,
		constants.Capacity12)
	
	results = append(results, AnyOs)
	
	if osconsts.IsWindows {
		results = append(results, Windows)
	} else {
		results = append(results, Unix)
	}
	
	if osconsts.IsLinux {
		results = append(results, Linux)
	}
	
	if osconsts.IsDarwinOrMacOs {
		results = append(results, MacOs)
	}
	
	if osconsts.CurrentOperatingSystem == osconsts.Android {
		results = append(results, Android)
	}
	
	if osconsts.IsFreebsd {
		results = append(results, FreeBsd)
	}
	
	if IsRunningInDockerContainer() {
		results = append(results, Docker)
	}
	
	currentOsMixType := CurrentOsType()
	
	switch currentOsMixType {
	case Ubuntu:
		results = append(results, Ubuntu)
	case Centos:
		results = append(results, Centos)
	case Debian:
		results = append(results, Debian)
	case RedHatEnterpriseLinux:
		results = append(results, RedHatEnterpriseLinux)
	}
	
	return results
}

func (it generate) OperatingSystemDetailLazy() (detail *OperatingSystemDetail, err error) {
	existing, err := it.getOperatingSystemDetailUsingFs()
	isFirstTime := err == nil && existing == nil
	
	if isFirstTime {
		// generate and store
		return it.operatingSystemDetailGenerateSave()
	}
	
	isReadSuccessFromFileSystem :=
		err == nil && existing != nil
	
	if isReadSuccessFromFileSystem {
		return existing.OperatingSystemDetail, errcore.ToError(existing.Error)
	}
	
	isReadSuccessFromFileSystemWithErr :=
		err != nil && existing != nil
	
	if isReadSuccessFromFileSystemWithErr {
		// remove cache
		// warning intentionally:
		//  swallowing the error
		//  as it is cache and has issues so removing it
		_ = os.RemoveAll(osDetailTempCacheRootPath)
	}
	
	// make a fresh start
	return it.operatingSystemDetailGenerateSave()
}

func (it generate) getOperatingSystemDetailUsingFs() (*OsDetailWithErr, error) {
	isExist := chmodhelper.IsPathExists(
		osDetailTempCachePath)
	
	if !isExist {
		return nil, nil
	}
	
	// exist
	jsonResult := strtype.
		NewFileReader(osDetailTempCachePath).
		RawAsJsonResult()
	
	var detailWithErr OsDetailWithErr
	err := jsonResult.Deserialize(&detailWithErr)
	
	return &detailWithErr, err
}

func (it generate) saveOperatingSystemDetailUsingFs(
	detail *OsDetailWithErr,
) error {
	err := it.createTempDirOnRequired()
	
	if err != nil {
		return err
	}
	
	json := detail.Json()
	
	if json.HasIssuesOrEmpty() {
		return json.MeaningfulError()
	}
	
	writeErr := ioutil.WriteFile(
		osDetailTempCachePath,
		json.Bytes,
		cacheFileMode)
	
	if writeErr != nil {
		return writeErr
	}
	
	if osconsts.IsWindows {
		return nil
	}
	
	return chmodhelper.
		ChmodApply.
		OnMismatch(
			false,
			cacheFileMode,
			osDetailTempCachePath)
}

func (it generate) createTempDirOnRequired() error {
	if chmodhelper.IsPathExists(osDetailTempCacheRootPath) {
		return nil
	}
	
	err := os.MkdirAll(
		osDetailTempCacheRootPath,
		filemode.DirDefault)
	
	if err != nil {
		return err
	}
	
	return nil
}

func (it generate) OperatingSystemDetail() (detail *OperatingSystemDetail, err error) {
	if osconsts.IsWindows {
		return it.windowsOperatingSystemDetail()
	}
	
	// unix
	return it.unixOperatingSystemDetail()
}

func (it generate) unixOperatingSystemDetail() (*OperatingSystemDetail, error) {
	if osconsts.IsDarwinOrMacOs {
		return it.macOsOperatingSystemDetail()
	}
	
	// other linux
	return it.linuxOperatingSystemDetail()
}

func (it generate) macOsOperatingSystemDetail() (*OperatingSystemDetail, error) {
	cmd := exec.Command(macOsSysVersionCommand)
	compiledOutput, err := cmd.CombinedOutput()
	
	if err != nil {
		return nil, errcore.FailedToExecuteType.ErrorRefOnly(macOsSysVersionCommand)
	}
	
	return it.ProcessMacOsOutputLines(compiledOutput)
}

// windowsOperatingSystemDetail
//
//	Generates Operating System Details for Windows
//
// References:
//   - Our issue : https://github.com/alimtvnetwork/enum-v8/-/issues/4
func (it generate) windowsOperatingSystemDetail() (*OperatingSystemDetail, error) {
	sysDetailGetter, err := getWinSysDetail()
	
	if err != nil {
		return nil, err
	}
	
	return sysDetailGetter.SystemDetail()
}

// ProcessMacOsOutputLines
//
// OutputLines:
//
//	ProductName:	Mac OS X
//	ProductVersion:	10.15.7
//	BuildVersion:	19H524
func (it generate) ProcessMacOsOutputLines(
	outputLines []byte,
) (*OperatingSystemDetail, error) {
	toString := string(outputLines)
	
	if toString == "" {
		return nil, errcore.NotSupportedType.Error(
			"empty outputLines-lines",
			macOsSysVersionCommand)
	}
	
	colonOutputLinesMap := it.keyValuesColonLinesToMap(
		strings.Split(toString, constants.DefaultLine),
	)
	
	if len(colonOutputLinesMap) == 0 {
		return nil, errcore.NotSupportedType.Error(
			"couldn't be able to process mac version!\n"+toString,
			macOsSysVersionCommand)
	}
	
	productName := colonOutputLinesMap[macOsProductName]       // eg. Mac OS X, "ProductName"
	productVersion := colonOutputLinesMap[macOsProductVersion] // eg. 10.15.7, "ProductVersion"
	buildVersion := colonOutputLinesMap[macOsBuildVersion]     // eg. 19H524, "BuildVersion"
	
	finalName := strtype.New(productName)
	
	return &OperatingSystemDetail{
		Name:         finalName, // eg, Mac OS X
		ProductName:  finalName,
		Vendor:       macOsName,
		Version:      strtype.New(productVersion), // eg. 10.15.7
		Release:      strtype.New(buildVersion),   // build version, eg. 19H524
		Architecture: osarchs.CurrentArch,
		OsMixType:    MacOs,
		Group:        ostype.CurrentGroup,
		IsMacOs:      true,
		IsDocker:     IsRunningInDockerContainer(),
	}, nil
}

// keyValuesColonLinesToMap
//
// each line contains key : {whitespace} {value}
func (it generate) keyValuesColonLinesToMap(colonSeparatorLines []string) map[string]string {
	if len(colonSeparatorLines) == 0 {
		return map[string]string{}
	}
	
	newMap := make(map[string]string, len(colonSeparatorLines)+1)
	
	for _, line := range colonSeparatorLines {
		trimmedLine := strings.TrimSpace(line)
		
		if trimmedLine == "" {
			continue
		}
		
		left, right := stringutil.SplitLeftRightTrimmed(
			line,
			constants.Colon)
		
		newMap[left] = right
	}
	
	return newMap
}

// linuxOperatingSystemDetail
//
//	Generates Operating System Details for linux
//
// References:
//   - SysInfo Package                 : https://github.com/zcalusic/sysinfo
//   - SysInfo Package (Specific file) : https://github.com/zcalusic/sysinfo/blob/master/os.go
func (it generate) linuxOperatingSystemDetail() (*OperatingSystemDetail, error) {
	defaultLinuxReleaseFile, err := strtype.NewFileReader(
		linuxvendortype.DefaultLinuxReleasePath,
	).OpenFile()
	
	if err != nil {
		return nil, err
	}
	
	defer func() { _ = defaultLinuxReleaseFile.Close() }()
	
	var (
		name                   string
		vendor                 string
		version                string
		release                string
		osMixType              = Linux
		vendorType             = linuxvendortype.Invalid
		prettyNameRegex        = prettyNameLazyRegex.CompileMust()        // once
		exactIdFieldMatchRegex = exactIdFieldMatchLazyRegex.CompileMust() // once
		versionIdRegex         = versionIdLazyRegex.CompileMust()         // once
		ubuntuRegex            = ubuntuLazyRegex.CompileMust()            // once
		centOSRegex            = centOSLazyRegex.CompileMust()            // once
		redHatRegex            = redHatLazyRegex.CompileMust()            // once
	)
	
	s := bufio.NewScanner(defaultLinuxReleaseFile)
	for s.Scan() {
		if m := prettyNameRegex.FindStringSubmatch(s.Text()); m != nil {
			name = strings.Trim(m[1], `"`)
		} else if m := exactIdFieldMatchRegex.FindStringSubmatch(s.Text()); m != nil {
			vendor = strings.Trim(m[1], `"`)
		} else if m := versionIdRegex.FindStringSubmatch(s.Text()); m != nil {
			version = strings.Trim(m[1], `"`)
		}
	}
	
	switch vendor {
	case linuxvendortype.Debian.ComparingName():
		osMixType = Debian
		vendorType = linuxvendortype.Debian
		release = readTrimmedFile(vendorType.ReleaseInfoFilePath())
	case linuxvendortype.Ubuntu.ComparingName():
		osMixType = Ubuntu
		vendorType = linuxvendortype.Ubuntu
		
		if m := ubuntuRegex.FindStringSubmatch(name); m != nil {
			release = m[1]
		}
	case linuxvendortype.CentOs.ComparingName():
		osMixType = Centos
		vendorType = linuxvendortype.CentOs
		
		if rel := readTrimmedFile(vendorType.ReleaseInfoFilePath()); rel != "" {
			if m := centOSRegex.FindStringSubmatch(rel); m != nil {
				release = m[2]
			}
		}
	case linuxvendortype.RHEL.ComparingName():
		osMixType = RedHatEnterpriseLinux
		vendorType = linuxvendortype.RHEL
		if rel := readTrimmedFile(vendorType.ReleaseInfoFilePath()); rel != "" {
			if m := redHatRegex.FindStringSubmatch(rel); m != nil {
				release = m[1]
			}
		}
		
		if release == "" {
			if m := redHatRegex.FindStringSubmatch(name); m != nil {
				release = m[1]
			}
		}
	}
	
	finalName := strtype.New(name)
	
	return &OperatingSystemDetail{
		Name:            finalName,
		ProductName:     finalName,
		Vendor:          strtype.New(version),
		LinuxVendorType: vendorType,
		Version:         strtype.New(version),
		Release:         strtype.New(release),
		Architecture:    osarchs.CurrentArch,
		OsMixType:       osMixType,
		Group:           ostype.CurrentGroup,
		IsLinux:         true,
		IsDocker:        IsRunningInDockerContainer(),
	}, nil
}

func (it generate) operatingSystemDetailGenerateSave() (*OperatingSystemDetail, error) {
	osDetail, err := it.OperatingSystemDetail()
	
	// swallowed error intentionally:
	//  as if the cache couldn't save doesn't matter
	//  regenerate should be fine.
	_ = it.saveOperatingSystemDetailUsingFs(&OsDetailWithErr{
		OperatingSystemDetail: osDetail,
		Error:                 errcore.ToString(err),
	})
	
	return osDetail, err
}

func (it generate) currentOsMixTypesMap() map[Variant]bool {
	allItems := CurrentOsMixTypes()
	resultMap := make(map[Variant]bool, len(allItems)+1)
	
	for _, item := range allItems {
		resultMap[item] = true
	}
	
	return resultMap
}
