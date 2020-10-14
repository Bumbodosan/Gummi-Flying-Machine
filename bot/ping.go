package bot

import (
	"time"

	"github.com/bwmarrin/discordgo"
)

type PingCommand struct{}

func (c PingCommand) Run(bot *Bot, args string, message *discordgo.Message) Sendable {
	return Message{Content: "Pong!", Timeout: time.Second * 10}
}
