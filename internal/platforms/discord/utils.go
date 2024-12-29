package discord

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

func newStatus(name string, value any) discordgo.UpdateStatusData {
	return discordgo.UpdateStatusData{
		Status: "online",
		Activities: []*discordgo.Activity{
			{
				Type:     discordgo.ActivityTypeCustom,
				Name:     fmt.Sprintf("%s: %v", name, value),
				URL:      "",
				State:    fmt.Sprintf("%s: %v", name, value),
				Details:  fmt.Sprintf("%s: %v", name, value),
				Instance: true,
			},
		},
	}
}
