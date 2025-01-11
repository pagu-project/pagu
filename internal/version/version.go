package pagu

import "fmt"

type Version struct {
	Major uint8
	Minor uint8
	Patch uint8
}

var version = Version{
	Major: 0,
	Minor: 0,
	Patch: 5,
}

func StringVersion() string {
	ver := fmt.Sprintf("%d.%d.%d", version.Major, version.Minor, version.Patch)

	return ver
}
