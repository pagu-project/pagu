package entity

type PlatformID int

const (
	PlatformIDCLI      PlatformID = 1
	PlatformIDDiscord  PlatformID = 2
	PlatformIDWeb      PlatformID = 3
	PlatformIDReserved PlatformID = 4
	PlatformIDTelegram PlatformID = 5
)

func (appID PlatformID) String() string {
	switch appID {
	case PlatformIDCLI:
		return "CLI"
	case PlatformIDDiscord:
		return "Discord"
	case PlatformIDWeb:
		return "Web"
	case PlatformIDReserved:
		return "Reserved"
	case PlatformIDTelegram:
		return "Telegram"
	}

	return ""
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
