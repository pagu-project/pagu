package pagu

import "fmt"

type Version struct {
	Meta  string
	Major uint8
	Minor uint8
	Patch uint8
}

var version = Version{
	Major: 0,
	Minor: 0,
	Patch: 1,
	Meta:  "beta",
}

func StringVersion() string {
	ver := fmt.Sprintf("%d.%d.%d", version.Major, version.Minor, version.Patch)

	if version.Meta != "" {
		ver = fmt.Sprintf("%s-%s", ver, version.Meta)
	}

	return ver
}
