package bot

import (
	"time"

	"github.com/bwmarrin/discordgo"
)

type PingCommand struct{}

func (c PingCommand) Run(bot *Bot, args string, message *discordgo.Message) error {
	sentMessage, err := bot.Session.ChannelMessageSend(message.ChannelID, "Pong!")

	if err != nil {
		return err
	}

	time.Sleep(3 * time.Second)
	bot.Session.ChannelMessagesBulkDelete(message.ChannelID, []string{sentMessage.ID, message.ID})

	return nil
}
