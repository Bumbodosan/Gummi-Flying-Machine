package bot

import (
	"encoding/base64"
	"time"

	"github.com/bwmarrin/discordgo"
)

type EncodeBase64Command struct{}

func (c EncodeBase64Command) Run(bot *Bot, args string, message *discordgo.Message) error {
	s := base64.StdEncoding.EncodeToString([]byte(args))

	sentMessage, err := bot.Session.ChannelMessageSend(message.ChannelID, s)

	if err != nil {
		return err
	}

	time.Sleep(5 * time.Second)
	bot.Session.ChannelMessagesBulkDelete(message.ChannelID, []string{sentMessage.ID, message.ID})

	return nil
}
