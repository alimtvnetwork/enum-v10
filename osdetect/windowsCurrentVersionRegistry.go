package osdetect

// windowsCurrentVersionRegistry
//
// https://prnt.sc/Yf7QKMuJNTse
//
// Issue : https://github.com/alimtvnetwork/enum-v10/-/issues/4
type windowsCurrentVersionRegistry struct {
	buildBranch          string
	productName          string
	installationType     string // eg. Client (Win 10, 11) or Server (2016, 2019)
	editionId            string
	compositionEditionId string
	releaseId            string // eg. 2009, 1809
	pathName             string // represents c:\windows
	systemRoot           string // represents c:\windows
	currentVersion       string
	majorVersionNumber   string
	minorVersionNumber   string
	currentBuild         string
	currentBuildNumber   string
	registeredOwner      string
}
