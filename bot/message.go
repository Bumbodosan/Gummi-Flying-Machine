package bot

import (
	"io"
	"time"

	"github.com/bwmarrin/discordgo"
)

type Sendable interface {
	Send(*Bot, *discordgo.Message) error
}

type Message struct {
	Content string
	Files   []File
	Timeout time.Duration // zero for don't remove
	// maybe different timeout for message being replied to

	// Add more things from discordgo.MessageSend here when needed
}

type File struct {
	Name   string
	Reader io.Reader
}

type ErrorMessage struct {
	Content string

	// Add more things from discordgo.MessageSend here when needed
}

func (msg Message) Send(bot *Bot, replyingTo *discordgo.Message) error {
	var msgsSent []*discordgo.Message
	if msg.Content != "" {
		ms, err := bot.Session.ChannelMessageSend(replyingTo.ChannelID, msg.Content)
		if err != nil {
			return err
		}
		msgsSent = append(msgsSent, ms)
	}

	if msg.Files != nil {
		for _, file := range msg.Files {
			ms, err := bot.Session.ChannelMessageSendComplex(
				replyingTo.ChannelID,
				&discordgo.MessageSend{
					Files: []*discordgo.File{
						{
							Name:   file.Name,
							Reader: file.Reader,
						},
					},
				},
			)
			if err != nil {
				return err
			}
			msgsSent = append(msgsSent, ms)
		}
	}

	if msg.Timeout != 0 {
		time.Sleep(msg.Timeout)

		for _, msgSent := range msgsSent {
			err := bot.Session.ChannelMessagesBulkDelete(
				replyingTo.ChannelID,
				[]string{replyingTo.ID, msgSent.ID},
			)

			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (msg ErrorMessage) Send(bot *Bot, replyingTo *discordgo.Message) error {
	msgSent, err := bot.Session.ChannelMessageSend(replyingTo.ChannelID, msg.Content)

	if err != nil {
		return err
	}

	time.Sleep(time.Second * 20)

	err = bot.Session.ChannelMessagesBulkDelete(
		replyingTo.ChannelID,
		[]string{replyingTo.ID, msgSent.ID},
	)

	if err != nil {
		return err
	}

	return nil
}
