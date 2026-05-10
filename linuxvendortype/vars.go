package linuxvendortype

import (
	"github.com/alimtvnetwork/core-v9/coreimpl/enumimpl"
)

var (
	Ranges = [...]string{
		Invalid:   "Invalid",
		Ubuntu:    "Ubuntu",
		Debian:    "Debian",
		LinuxMint: "LinuxMint",
		CentOs:    "CentOs",
		RHEL:      "RHEL",
		Gentoo:    "Gentoo",
		Fedora:    "Fedora",
		Kali:      "Kali",
		ArchLinux: "ArchLinux",
		OpenSuse:  "OpenSuse",
	}

	// github.com/alimtvnetwork/enum-v9/-/issues/4
	// https://t.ly/2KHe, https://prnt.sc/M9SPHl4GBYFN
	aliasMap = map[string]byte{
		"ubuntu": Ubuntu.ValueByte(),
		"debian": Debian.ValueByte(),
		"centos": CentOs.ValueByte(),
		"rhel":   RHEL.ValueByte(),
	}

	comparingNamesMap = [...]string{
		Ubuntu: "ubuntu",
		Debian: "debian",
		CentOs: "centos",
		RHEL:   "rhel",
	}

	releaseInfoFilePathMap = map[Variant]string{
		Debian: "/etc/debian_version", // https://t.ly/ZNY9
		CentOs: "/etc/centos-release",
		RHEL:   "/etc/redhat-release",
	}

	displayMap = map[Variant]string{
		Invalid:   "",
		Ubuntu:    "Ubuntu",
		Debian:    "Debian",
		LinuxMint: "Linux Mint",
		CentOs:    "Cent Os",
		RHEL:      "Red Hat Enterprise Linux", // https://www.tecmint.com/linux-distro-for-power-users/
		Gentoo:    "Gentoo",
		Fedora:    "Fedora",
		Kali:      "Kali Linux",
		ArchLinux: "ArchLinux Linux",
		OpenSuse:  "OpenSUSE",
	}

	BasicEnumImpl = enumimpl.New.BasicByte.UsingFirstItemSliceAliasMap(
		Invalid,
		Ranges[:],
		aliasMap)
)
