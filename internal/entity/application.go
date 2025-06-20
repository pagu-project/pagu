package entity

import (
	"fmt"

	"gopkg.in/yaml.v3"
)

// BotID defines the identifier used to initialize a bot instance.
// It is primarily used to filter out commands not intended for a specific bot.
// For example, moderator commands are only available to the Moderator bot.
type BotID int

const (
	BotID_CLI       BotID = 1 //nolint // underscores used for BotID
	BotID_Discord   BotID = 2 //nolint // underscores used for BotID
	BotID_Moderator BotID = 3 //nolint // underscores used for BotID
	BotID_Telegram  BotID = 4 //nolint // underscores used for BotID
	BotID_WhatsApp  BotID = 5 //nolint // underscores used for BotID
	BotID_Web       BotID = 6 //nolint // underscores used for BotID
)

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

func (b BotID) String() string {
	for name, id := range BotNameToID {
		if id == b {
			return name
		}
	}

	return fmt.Sprintf("%d", b)
}

// PlatformID defines the platform from which the user is calling the API.
// It is stored in the database alongside the user ID to track user activity on specific platforms.
// The numeric values must be preserved for consistency.
type PlatformID int

const (
	PlatformIDCLI      PlatformID = 1
	PlatformIDDiscord  PlatformID = 2
	PlatformIDWeb      PlatformID = 3
	PlatformIDWhatsapp PlatformID = 4
	PlatformIDTelegram PlatformID = 5
)

var platformIDToString = map[PlatformID]string{
	PlatformIDCLI:      "CLI",
	PlatformIDDiscord:  "Discord",
	PlatformIDWeb:      "Web",
	PlatformIDWhatsapp: "Whatsapp",
	PlatformIDTelegram: "Telegram",
}

func (pid PlatformID) String() string {
	str, ok := platformIDToString[pid]
	if ok {
		return str
	}

	return fmt.Sprintf("%d", pid)
}
