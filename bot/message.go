package bot

import (
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
)

type Sendable interface {
	Send(*Bot, *discordgo.Message)
}

type Message struct {
	Content string
	Timeout time.Duration // zero for don't remove
	// maybe different timeout for message being replied to

	// Add more things from discordgo.MessageSend here when needed
}

type ErrorMessage struct {
	Content string

	// Add more things from discordgo.MessageSend here when needed
}

func (msg Message) Send(bot *Bot, replyingTo *discordgo.Message) {
	msgSent, err := bot.Session.ChannelMessageSend(replyingTo.ChannelID, msg.Content)

	if err != nil {
		fmt.Println("discordgo error: ", err)
		return
	}

	if msg.Timeout != 0 {
		time.Sleep(msg.Timeout)

		err = bot.Session.ChannelMessagesBulkDelete(
			replyingTo.ChannelID,
			[]string{replyingTo.ID, msgSent.ID},
		)

		if err != nil {
			fmt.Println("discordgo error: ", err)
			return
		}
	}
}

func (msg ErrorMessage) Send(bot *Bot, replyingTo *discordgo.Message) {
	msgSent, err := bot.Session.ChannelMessageSend(replyingTo.ChannelID, msg.Content)

	if err != nil {
		fmt.Println("discordgo error: ", err)
		return
	}

	time.Sleep(time.Second * 20)

	err = bot.Session.ChannelMessagesBulkDelete(
		replyingTo.ChannelID,
		[]string{replyingTo.ID, msgSent.ID},
	)

	if err != nil {
		fmt.Println("discordgo error: ", err)
		return
	}
}
