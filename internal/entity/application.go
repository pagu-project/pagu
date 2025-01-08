package entity

import (
	"fmt"
)

type PlatformID int

const (
	PlatformIDCLI      PlatformID = 1
	PlatformIDDiscord  PlatformID = 2
	PlatformIDWeb      PlatformID = 3
	PlatformIDReserved PlatformID = 4
	PlatformIDTelegram PlatformID = 5
)

var platformIDToString = map[PlatformID]string{
	PlatformIDCLI:      "CLI",
	PlatformIDDiscord:  "Discord",
	PlatformIDWeb:      "Web",
	PlatformIDReserved: "Reserved",
	PlatformIDTelegram: "Telegram",
}

func (pid PlatformID) String() string {
	str, ok := platformIDToString[pid]
	if ok {
		return str
	}

	return fmt.Sprintf("%d", pid)
}

func AllAppIDs() []PlatformID {
	return []PlatformID{
		PlatformIDCLI,
		PlatformIDDiscord,
		PlatformIDWeb,
		PlatformIDReserved,
		PlatformIDTelegram,
	}
}
