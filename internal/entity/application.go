package entity

import (
	"fmt"

	"gopkg.in/yaml.v3"
)

type BotID int

const (
	BotID_CLI       BotID = 1 //nolint // underscores used for BotID
	BotID_Discord   BotID = 2 //nolint // underscores used for BotID
	BotID_Moderator BotID = 3 //nolint // underscores used for BotID
	BotID_Telegram  BotID = 4 //nolint // underscores used for BotID
	BotID_WhatsApp  BotID = 5 //nolint // underscores used for BotID
	BotID_Web       BotID = 6 //nolint // underscores used for BotID
)

var BotNameToID = map[string]BotID{
	"CLI":       BotID_CLI,
	"Discord":   BotID_Discord,
	"Moderator": BotID_Moderator,
	"Telegram":  BotID_Telegram,
	"WhatsApp":  BotID_WhatsApp,
	"Web":       BotID_Web,
}

func (b *BotID) UnmarshalYAML(value *yaml.Node) error {
	var botName string
	if err := value.Decode(&botName); err != nil {
		return err
	}

	id, exists := BotNameToID[botName]
	if !exists {
		return fmt.Errorf("invalid bot name: %s", botName)
	}

	*b = id

	return nil
}

type PlatformID int

const (
	PlatformIDCLI      PlatformID = 1
	PlatformIDDiscord  PlatformID = 2
	PlatformIDWeb      PlatformID = 3
	PlatformIDReserved PlatformID = 4
	PlatformIDTelegram PlatformID = 5
	PlatformIDWhatsapp PlatformID = 6
)

var platformIDToString = map[PlatformID]string{
	PlatformIDCLI:      "CLI",
	PlatformIDDiscord:  "Discord",
	PlatformIDWeb:      "Web",
	PlatformIDReserved: "Reserved",
	PlatformIDTelegram: "Telegram",
	PlatformIDWhatsapp: "Whatsapp",
}

func (pid PlatformID) String() string {
	str, ok := platformIDToString[pid]
	if ok {
		return str
	}

	return fmt.Sprintf("%d", pid)
}

func AllBotIDs() []BotID {
	return []BotID{
		BotID_CLI,
		BotID_Discord,
		BotID_Moderator,
		BotID_Telegram,
		BotID_WhatsApp,
		BotID_Web,
	}
}
